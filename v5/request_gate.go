package okx

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RequestGateConfig 用于配置 REST 请求闸门（并发 + 速率）。
//
// 说明：
// - 目标是“前置流控”，避免高并发下触发 429 风暴。
// - 这是客户端级别的闸门（每个 Client 一份），不会跨进程/跨实例共享配额。
type RequestGateConfig struct {
	// MaxConcurrent 限制同时在途的 HTTP 请求数。
	// <= 0 表示不限制。
	MaxConcurrent int

	// GlobalRPS 表示全局每秒请求数限制（token bucket）。
	// <= 0 表示不限制。
	GlobalRPS float64

	// GlobalBurst 表示全局突发容量（token bucket 容量上限）。
	// <= 0 时会使用一个安全默认值（当 GlobalRPS>0 时默认 1）。
	GlobalBurst int
}

func defaultRequestGateConfig() RequestGateConfig {
	return RequestGateConfig{
		// 默认做“并发 + 保守速率”双闸门：
		// - 并发用于削峰，避免单机突发把 OKX 打进限频/网关异常；
		// - 全局 RPS 只是安全下限（不同账户/接口限频差异很大），生产环境应按自身场景显式配置或禁用。
		MaxConcurrent: 10,
		GlobalRPS:     10,
		GlobalBurst:   20,
	}
}

// WithRequestGate 设置 REST 请求闸门（并发 + 速率）。
func WithRequestGate(cfg RequestGateConfig) Option {
	cfgCopy := cfg
	if cfgCopy.MaxConcurrent < 0 {
		cfgCopy.MaxConcurrent = 0
	}
	if cfgCopy.GlobalRPS < 0 {
		cfgCopy.GlobalRPS = 0
	}
	if cfgCopy.GlobalBurst < 0 {
		cfgCopy.GlobalBurst = 0
	}

	return func(c *Client) {
		c.gate = newRequestGate(cfgCopy)
	}
}

// WithRequestGateDisabled 关闭 REST 请求闸门。
func WithRequestGateDisabled() Option {
	return func(c *Client) {
		c.gate = nil
	}
}

type routeKey struct {
	Method   string
	Endpoint string
}

const requestGateMethodWS = "WS"

func wsOpGateKey(op string) string {
	return "op:" + op
}

type requestGate struct {
	sem *semaphore

	globalLimiter *tokenBucketLimiter

	mu           sync.RWMutex
	routeLimiter map[routeKey]*tokenBucketLimiter
}

func newRequestGate(cfg RequestGateConfig) *requestGate {
	g := &requestGate{
		routeLimiter: make(map[routeKey]*tokenBucketLimiter),
	}
	if cfg.MaxConcurrent > 0 {
		g.sem = newSemaphore(cfg.MaxConcurrent)
	}
	if cfg.GlobalRPS > 0 {
		burst := cfg.GlobalBurst
		if burst <= 0 {
			burst = 1
		}
		g.globalLimiter = newTokenBucketLimiter(cfg.GlobalRPS, burst)
	}
	return g
}

func (g *requestGate) acquire(ctx context.Context, method, endpoint string) (release func(), err error) {
	if g == nil {
		return func() {}, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	release = func() {}
	if g.sem != nil {
		if err := g.sem.Acquire(ctx); err != nil {
			return nil, err
		}
		release = g.sem.Release
	}

	if g.globalLimiter != nil {
		if err := g.globalLimiter.Wait(ctx); err != nil {
			release()
			return nil, err
		}
	}

	var rl *tokenBucketLimiter
	if method != "" && endpoint != "" {
		g.mu.RLock()
		rl = g.routeLimiter[routeKey{Method: method, Endpoint: endpoint}]
		g.mu.RUnlock()
	}
	if rl != nil {
		if err := rl.Wait(ctx); err != nil {
			release()
			return nil, err
		}
	}

	return release, nil
}

func (g *requestGate) setRouteLimiter(keys []routeKey, limiter *tokenBucketLimiter) {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, k := range keys {
		if limiter == nil {
			delete(g.routeLimiter, k)
			continue
		}
		g.routeLimiter[k] = limiter
	}
}

func (c *Client) applyTradeAccountRateLimit(info *TradeAccountRateLimit) error {
	if c == nil || c.gate == nil || info == nil {
		return nil
	}

	// OKX 文档：accRateLimit 表示“2 秒内可用次数”（Place/Amend，含 REST/WS）。
	limit2s, err := strconv.ParseInt(info.AccRateLimit, 10, 64)
	if err != nil {
		return fmt.Errorf("okx: invalid accRateLimit: %w", err)
	}
	effective2s := limit2s
	if effective2s <= 0 && info.NextAccRateLimit != "" {
		if next2s, err := strconv.ParseInt(info.NextAccRateLimit, 10, 64); err == nil && next2s > 0 {
			effective2s = next2s
		}
	}
	if effective2s <= 0 {
		// Fail-Safe：口径不明确或额度为 0 时，至少保留一个极保守的闸门，避免“Fail-Open”放大限频风暴。
		effective2s = 1
	}

	rps := float64(effective2s) / 2.0
	burst := int(effective2s)
	limiter := newTokenBucketLimiter(rps, burst)
	c.gate.setRouteLimiter(tradeAccountRateLimitKeys(), limiter)
	return nil
}

func tradeAccountRateLimitKeys() []routeKey {
	return append(tradeAccountRateLimitRESTKeys(), tradeAccountRateLimitWSKeys()...)
}

func tradeAccountRateLimitRESTKeys() []routeKey {
	return []routeKey{
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/order"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/batch-orders"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/amend-order"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/amend-batch-orders"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/order-algo"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/amend-algos"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/cancel-order"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/cancel-batch-orders"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/cancel-algos"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/mass-cancel"},
		{Method: http.MethodPost, Endpoint: "/api/v5/trade/cancel-all-after"},
	}
}

func tradeAccountRateLimitWSKeys() []routeKey {
	return []routeKey{
		{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpOrder)},
		{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpCancelOrder)},
		{Method: requestGateMethodWS, Endpoint: wsOpGateKey(wsOpAmendOrder)},
	}
}

type semaphore struct {
	ch chan struct{}
}

func newSemaphore(max int) *semaphore {
	return &semaphore{ch: make(chan struct{}, max)}
}

func (s *semaphore) Acquire(ctx context.Context) error {
	if s == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *semaphore) Release() {
	if s == nil {
		return
	}
	select {
	case <-s.ch:
	default:
	}
}

// tokenBucketLimiter 是一个轻量的 token bucket（无外部依赖，支持 ctx.Wait）。
// 仅用于请求闸门，不承诺强一致精度；目标是“削峰 + 避免紧循环”。
type tokenBucketLimiter struct {
	mu sync.Mutex

	rps   float64
	burst float64

	tokens float64
	last   time.Time
}

func newTokenBucketLimiter(rps float64, burst int) *tokenBucketLimiter {
	if rps <= 0 || burst <= 0 {
		return nil
	}
	now := time.Now()
	b := float64(burst)
	return &tokenBucketLimiter{
		rps:    rps,
		burst:  b,
		tokens: b,
		last:   now,
	}
}

func (l *tokenBucketLimiter) Wait(ctx context.Context) error {
	if l == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	for {
		wait := l.takeOrComputeWait()
		if wait <= 0 {
			return nil
		}

		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
			// 继续循环补充 token
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
}

func (l *tokenBucketLimiter) takeOrComputeWait() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.rps <= 0 || l.burst <= 0 {
		return 0
	}

	now := time.Now()
	if !l.last.IsZero() {
		elapsed := now.Sub(l.last).Seconds()
		if elapsed > 0 {
			l.tokens += elapsed * l.rps
			if l.tokens > l.burst {
				l.tokens = l.burst
			}
		}
	}
	l.last = now

	if l.tokens >= 1 {
		l.tokens -= 1
		return 0
	}

	need := 1 - l.tokens
	seconds := need / l.rps
	if seconds <= 0 {
		return 0
	}

	d := time.Duration(seconds * float64(time.Second))
	if d <= 0 {
		return time.Millisecond
	}
	return d
}
