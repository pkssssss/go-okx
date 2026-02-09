package okx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkssssss/go-okx/v5/internal/rest"
	"github.com/pkssssss/go-okx/v5/internal/sign"
)

const (
	defaultBaseURL   = "https://www.okx.com"
	defaultUserAgent = "OKX/go-okx"
)

// Credentials 是 OKX APIKey 的三元组。
type Credentials struct {
	APIKey     string
	SecretKey  string
	Passphrase string
}

func (c Credentials) Redacted() Credentials {
	return Credentials{
		APIKey:     maskLast4(c.APIKey),
		SecretKey:  "***",
		Passphrase: "***",
	}
}

func (c Credentials) String() string {
	r := c.Redacted()
	return fmt.Sprintf("Credentials{APIKey:%q, SecretKey:%q, Passphrase:%q}", r.APIKey, r.SecretKey, r.Passphrase)
}

func (c Credentials) GoString() string { return c.String() }

func (c Credentials) MarshalJSON() ([]byte, error) {
	r := c.Redacted()
	type out struct {
		APIKey     string `json:"APIKey"`
		SecretKey  string `json:"SecretKey"`
		Passphrase string `json:"Passphrase"`
	}
	return json.Marshal(out{
		APIKey:     r.APIKey,
		SecretKey:  r.SecretKey,
		Passphrase: r.Passphrase,
	})
}

func maskLast4(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "***"
	}
	return "****" + s[len(s)-4:]
}

// ClientErrorHandler 用于接收客户端内部“非业务请求”的运行错误（例如自动预热/后台任务）。
// 注意：正常 API 调用的错误仍通过 Do() 返回值返回，不会走该 handler。
type ClientErrorHandler func(err error)

// Client 是 OKX V5 API 客户端（REST + WS）。
// v0.1 先落地 REST 核心能力，WS 随后补齐。
type Client struct {
	rest *rest.Client

	creds *Credentials
	demo  bool

	gate *requestGate

	retry *RetryConfig

	errHandler ClientErrorHandler

	tradeAccountRateLimitMu          sync.Mutex
	tradeAccountRateLimitPrimed      atomic.Bool
	tradeAccountRateLimitLastAttempt atomic.Int64
	tradeAccountRateLimitLastErr     atomic.Value

	timeOffsetNanos atomic.Int64
	now             func() time.Time
}

// Option 用于配置 Client。
type Option func(*Client)

// NewClient 创建 OKX Client。
func NewClient(opts ...Option) *Client {
	c := &Client{
		rest: &rest.Client{
			BaseURL:   defaultBaseURL,
			UserAgent: defaultUserAgent,
		},
		gate: newRequestGate(defaultRequestGateConfig()),
		now:  time.Now,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithCredentials 设置 APIKey/Secret/Passphrase。
func WithCredentials(creds Credentials) Option {
	return func(c *Client) {
		c.creds = &creds
	}
}

// WithDemoTrading 开启/关闭模拟盘（REST 会附加 x-simulated-trading: 1，WS 使用 wspap host）。
func WithDemoTrading(enable bool) Option {
	return func(c *Client) {
		c.demo = enable
	}
}

// WithHTTPClient 设置自定义 HTTPClient。
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.rest.HTTPClient = hc
	}
}

// WithUserAgent 设置 User-Agent。
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.rest.UserAgent = userAgent
	}
}

// WithBaseURL 设置 REST BaseURL（默认 https://www.okx.com）。
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.rest.BaseURL = baseURL
	}
}

// WithTimeOffset 设置本地时间与服务器时间的偏移（用于签名时间戳校准）。
// 约定：签名时使用 now() - offset。
func WithTimeOffset(offset time.Duration) Option {
	return func(c *Client) {
		c.timeOffsetNanos.Store(offset.Nanoseconds())
	}
}

// WithNowFunc 覆盖当前时间函数（主要用于测试）。
func WithNowFunc(now func() time.Time) Option {
	return func(c *Client) {
		c.now = now
	}
}

// WithClientErrorHandler 设置客户端内部错误回调（用于可观测告警）。
func WithClientErrorHandler(handler ClientErrorHandler) Option {
	return func(c *Client) {
		c.errHandler = handler
	}
}

// RetryConfig 控制 REST 请求的重试策略（仅幂等 GET）。
//
// 注意：
// - 默认不启用重试（MaxRetries=0）。
// - 重试会重新生成签名时间戳（若为签名请求）。
type RetryConfig struct {
	// MaxRetries 表示最多重试次数（不含第一次请求）。
	MaxRetries int

	// BaseDelay 表示第一次重试的等待时间（指数退避起点）。
	// 未设置（<=0）且启用重试时，会使用安全默认值（200ms）。
	BaseDelay time.Duration

	// MaxDelay 表示最大等待时间（为 0 表示不限制）。
	MaxDelay time.Duration

	// RetryOnRateLimit 为 true 时，会在限速错误（HTTP 429 / code 50011/50061）上重试。
	RetryOnRateLimit bool
}

// WithRetry 设置重试策略（仅幂等 GET）。
func WithRetry(cfg RetryConfig) Option {
	cfgCopy := cfg
	if cfgCopy.MaxRetries < 0 {
		cfgCopy.MaxRetries = 0
	}
	if cfgCopy.BaseDelay < 0 {
		cfgCopy.BaseDelay = 0
	}
	if cfgCopy.MaxDelay < 0 {
		cfgCopy.MaxDelay = 0
	}
	if cfgCopy.MaxRetries > 0 && cfgCopy.BaseDelay <= 0 {
		cfgCopy.BaseDelay = 200 * time.Millisecond
	}

	return func(c *Client) {
		c.retry = &cfgCopy
	}
}

var errMissingCredentials = errors.New("okx: missing credentials")

type responseEnvelope struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (c *Client) do(ctx context.Context, method, endpoint string, query url.Values, body any, signed bool, out any) error {
	return c.doWithHeaders(ctx, method, endpoint, query, body, signed, nil, out)
}

func (c *Client) doWithHeaders(ctx context.Context, method, endpoint string, query url.Values, body any, signed bool, extraHeader http.Header, out any) error {
	requestPath := rest.BuildRequestPath(endpoint, query)

	var bodyBytes []byte
	var bodyString string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyBytes = b
		bodyString = string(b)
	}

	retryCfg := c.retry
	maxRetries := 0
	if retryCfg != nil && retryCfg.MaxRetries > 0 && method == http.MethodGet {
		maxRetries = retryCfg.MaxRetries
	}

	for attempt := 0; ; attempt++ {
		attemptCtx, attemptCancel := c.rest.ContextWithDefaultTimeout(ctx)

		if signed {
			if c.creds == nil || c.creds.APIKey == "" || c.creds.SecretKey == "" || c.creds.Passphrase == "" {
				if attemptCancel != nil {
					attemptCancel()
				}
				return errMissingCredentials
			}
		}

		if signed && isTradeAccountRateLimitedREST(method, endpoint) {
			if err := c.ensureTradeAccountRateLimit(attemptCtx); err != nil {
				if attemptCancel != nil {
					attemptCancel()
				}
				return &RequestStateError{
					Stage:       RequestStagePreflight,
					Dispatched:  false,
					Method:      method,
					RequestPath: requestPath,
					Err:         err,
				}
			}
		}

		release, err := c.gate.acquire(attemptCtx, method, endpoint)
		if err != nil {
			if attemptCancel != nil {
				attemptCancel()
			}
			return &RequestStateError{
				Stage:       RequestStageGate,
				Dispatched:  false,
				Method:      method,
				RequestPath: requestPath,
				Err:         err,
			}
		}

		header := make(http.Header)
		header.Set("Accept", "application/json")
		header.Set("Content-Type", "application/json")
		if c.demo {
			header.Set("x-simulated-trading", "1")
		}

		if signed {
			tm := c.now().Add(-c.TimeOffset())
			timestamp := sign.TimestampISO8601Millis(tm)
			prehash := sign.PrehashREST(timestamp, method, requestPath, bodyString)
			sig := sign.SignHMACSHA256Base64(c.creds.SecretKey, prehash)

			header.Set("OK-ACCESS-KEY", c.creds.APIKey)
			header.Set("OK-ACCESS-PASSPHRASE", c.creds.Passphrase)
			header.Set("OK-ACCESS-TIMESTAMP", timestamp)
			header.Set("OK-ACCESS-SIGN", sig)
		}

		if len(extraHeader) > 0 {
			for k, vs := range extraHeader {
				if len(vs) == 0 {
					continue
				}
				header.Del(k)
				for _, v := range vs {
					header.Add(k, v)
				}
			}
		}

		status, resp, respHeader, err := c.rest.Do(attemptCtx, method, requestPath, bodyBytes, header)
		release()
		if attemptCancel != nil {
			attemptCancel()
		}
		if err != nil {
			if attempt < maxRetries && isRetryableTransportError(err) {
				if err := sleepRetry(ctx, retryCfg, attempt+1); err != nil {
					return err
				}
				continue
			}
			return &RequestStateError{
				Stage:       RequestStageHTTP,
				Dispatched:  true,
				Method:      method,
				RequestPath: requestPath,
				Err:         err,
			}
		}

		if err := decodeEnvelope(status, resp, respHeader, method, requestPath, out); err != nil {
			if attempt < maxRetries && isRetryableAPIError(err, retryCfg) {
				if err := sleepRetry(ctx, retryCfg, attempt+1); err != nil {
					return err
				}
				continue
			}
			return err
		}
		return nil
	}
}

func (c *Client) doWithHeadersAndRequestID(ctx context.Context, method, endpoint string, query url.Values, body any, signed bool, extraHeader http.Header, out any) (requestID string, err error) {
	requestPath := rest.BuildRequestPath(endpoint, query)

	var bodyBytes []byte
	var bodyString string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		bodyBytes = b
		bodyString = string(b)
	}

	retryCfg := c.retry
	maxRetries := 0
	if retryCfg != nil && retryCfg.MaxRetries > 0 && method == http.MethodGet {
		maxRetries = retryCfg.MaxRetries
	}

	for attempt := 0; ; attempt++ {
		attemptCtx, attemptCancel := c.rest.ContextWithDefaultTimeout(ctx)

		if signed {
			if c.creds == nil || c.creds.APIKey == "" || c.creds.SecretKey == "" || c.creds.Passphrase == "" {
				if attemptCancel != nil {
					attemptCancel()
				}
				return "", errMissingCredentials
			}
		}

		if signed && isTradeAccountRateLimitedREST(method, endpoint) {
			if err := c.ensureTradeAccountRateLimit(attemptCtx); err != nil {
				if attemptCancel != nil {
					attemptCancel()
				}
				return requestID, &RequestStateError{
					Stage:       RequestStagePreflight,
					Dispatched:  false,
					Method:      method,
					RequestPath: requestPath,
					Err:         err,
				}
			}
		}

		release, err := c.gate.acquire(attemptCtx, method, endpoint)
		if err != nil {
			if attemptCancel != nil {
				attemptCancel()
			}
			return requestID, &RequestStateError{
				Stage:       RequestStageGate,
				Dispatched:  false,
				Method:      method,
				RequestPath: requestPath,
				Err:         err,
			}
		}

		header := make(http.Header)
		header.Set("Accept", "application/json")
		header.Set("Content-Type", "application/json")
		if c.demo {
			header.Set("x-simulated-trading", "1")
		}

		if signed {
			tm := c.now().Add(-c.TimeOffset())
			timestamp := sign.TimestampISO8601Millis(tm)
			prehash := sign.PrehashREST(timestamp, method, requestPath, bodyString)
			sig := sign.SignHMACSHA256Base64(c.creds.SecretKey, prehash)

			header.Set("OK-ACCESS-KEY", c.creds.APIKey)
			header.Set("OK-ACCESS-PASSPHRASE", c.creds.Passphrase)
			header.Set("OK-ACCESS-TIMESTAMP", timestamp)
			header.Set("OK-ACCESS-SIGN", sig)
		}

		if len(extraHeader) > 0 {
			for k, vs := range extraHeader {
				if len(vs) == 0 {
					continue
				}
				header.Del(k)
				for _, v := range vs {
					header.Add(k, v)
				}
			}
		}

		status, resp, respHeader, err := c.rest.Do(attemptCtx, method, requestPath, bodyBytes, header)
		release()
		if attemptCancel != nil {
			attemptCancel()
		}
		if respHeader != nil {
			requestID = respHeader.Get("x-request-id")
		}
		if err != nil {
			if attempt < maxRetries && isRetryableTransportError(err) {
				if err := sleepRetry(ctx, retryCfg, attempt+1); err != nil {
					return requestID, err
				}
				continue
			}
			return requestID, &RequestStateError{
				Stage:       RequestStageHTTP,
				Dispatched:  true,
				Method:      method,
				RequestPath: requestPath,
				Err:         err,
			}
		}

		if err := decodeEnvelope(status, resp, respHeader, method, requestPath, out); err != nil {
			if attempt < maxRetries && isRetryableAPIError(err, retryCfg) {
				if err := sleepRetry(ctx, retryCfg, attempt+1); err != nil {
					return requestID, err
				}
				continue
			}
			return requestID, err
		}
		return requestID, nil
	}
}

const tradeAccountRateLimitPrimeMinInterval = time.Second

// TradeAccountRateLimitPrimeStatus 表示 accRateLimit 自动预热的状态快照（用于自检/排障）。
type TradeAccountRateLimitPrimeStatus struct {
	Primed      bool
	LastAttempt time.Time
	LastError   string
}

func (c *Client) TradeAccountRateLimitPrimeStatus() TradeAccountRateLimitPrimeStatus {
	var s TradeAccountRateLimitPrimeStatus
	if c == nil {
		return s
	}

	s.Primed = c.tradeAccountRateLimitPrimed.Load()
	if ns := c.tradeAccountRateLimitLastAttempt.Load(); ns != 0 {
		s.LastAttempt = time.Unix(0, ns)
	}
	if v := c.tradeAccountRateLimitLastErr.Load(); v != nil {
		if msg, ok := v.(string); ok {
			s.LastError = msg
		}
	}
	return s
}

func (c *Client) onError(err error) {
	if c == nil || c.errHandler == nil || err == nil || errors.Is(err, context.Canceled) {
		return
	}
	defer func() {
		_ = recover()
	}()
	c.errHandler(err)
}

func (c *Client) ensureTradeAccountRateLimit(ctx context.Context) error {
	if c == nil || c.gate == nil || c.tradeAccountRateLimitPrimed.Load() {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	c.tradeAccountRateLimitMu.Lock()
	defer c.tradeAccountRateLimitMu.Unlock()

	if c.tradeAccountRateLimitPrimed.Load() {
		return nil
	}

	now := time.Now()
	if last := c.tradeAccountRateLimitLastAttempt.Load(); last != 0 {
		if now.Sub(time.Unix(0, last)) < tradeAccountRateLimitPrimeMinInterval {
			if v := c.tradeAccountRateLimitLastErr.Load(); v != nil {
				if msg, ok := v.(string); ok && msg != "" {
					return fmt.Errorf("okx: trade account-rate-limit prime failed recently: %s", msg)
				}
			}
			return nil
		}
	}
	c.tradeAccountRateLimitLastAttempt.Store(now.UnixNano())

	_, err := c.NewTradeAccountRateLimitService().Do(ctx)
	if err != nil {
		c.tradeAccountRateLimitLastErr.Store(err.Error())
		c.onError(fmt.Errorf("okx: trade account-rate-limit prime failed: %w", err))
		return err
	}

	c.tradeAccountRateLimitPrimed.Store(true)
	c.tradeAccountRateLimitLastErr.Store("")
	return nil
}

func isTradeAccountRateLimitedREST(method, endpoint string) bool {
	for _, k := range tradeAccountRateLimitRESTKeys() {
		if k.Method == method && k.Endpoint == endpoint {
			return true
		}
	}
	return false
}

func isRetryableTransportError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	return false
}

func isRetryableAPIError(err error, cfg *RetryConfig) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.HTTPStatus >= http.StatusInternalServerError {
		return true
	}

	if cfg != nil && cfg.RetryOnRateLimit {
		if apiErr.HTTPStatus == http.StatusTooManyRequests {
			return true
		}
		switch apiErr.Code {
		case "50011", "50061":
			return true
		}
	}
	return false
}

func sleepRetry(ctx context.Context, cfg *RetryConfig, retryIndex int) error {
	if cfg == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	d := retryDelay(*cfg, retryIndex)
	if d <= 0 {
		return nil
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func retryDelay(cfg RetryConfig, retryIndex int) time.Duration {
	if cfg.BaseDelay <= 0 || retryIndex <= 0 {
		return 0
	}

	d := cfg.BaseDelay
	for i := 1; i < retryIndex; i++ {
		if d > 0 && d > (1<<62) {
			break
		}
		d *= 2
	}

	if cfg.MaxDelay > 0 && d > cfg.MaxDelay {
		d = cfg.MaxDelay
	}

	// 加入抖动，避免多个 goroutine 同步重试造成“429 风暴”。
	// 采用 [d/2, d] 区间的均匀随机（Full Jitter 的简化变体）。
	if d <= 0 {
		return 0
	}
	half := d / 2
	if half <= 0 {
		return d
	}
	return half + retryJitterDuration(half)
}

var retryJitterRand = struct {
	mu  sync.Mutex
	rnd *rand.Rand
}{
	rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
}

func retryJitterDuration(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	retryJitterRand.mu.Lock()
	n := retryJitterRand.rnd.Int63n(int64(max) + 1)
	retryJitterRand.mu.Unlock()
	return time.Duration(n)
}

func decodeEnvelope(status int, resp []byte, respHeader http.Header, method, requestPath string, out any) error {
	var env responseEnvelope
	if err := json.Unmarshal(resp, &env); err != nil {
		if status < http.StatusBadRequest && out != nil {
			if err2 := json.Unmarshal(resp, out); err2 == nil {
				return nil
			}
		}
		return &APIError{
			HTTPStatus:  status,
			Method:      method,
			RequestPath: requestPath,
			Message:     "invalid JSON response",
			Raw:         resp,
			RequestID:   respHeader.Get("x-request-id"),
		}
	}

	// OKX 文档中少量接口会省略 code/msg/data 外层（直接返回对象/数组）。
	// 为提高 SDK 稳定性：在 HTTP 2xx 且缺少 envelope 字段时，尝试直接反序列化为 out。
	if env.Code == "" && env.Msg == "" && len(env.Data) == 0 {
		if status >= http.StatusBadRequest {
			return &APIError{
				HTTPStatus:  status,
				Method:      method,
				RequestPath: requestPath,
				Message:     "invalid response envelope",
				Raw:         resp,
				RequestID:   respHeader.Get("x-request-id"),
			}
		}
		if out == nil {
			return &APIError{
				HTTPStatus:  status,
				Method:      method,
				RequestPath: requestPath,
				Message:     "invalid response envelope",
				Raw:         resp,
				RequestID:   respHeader.Get("x-request-id"),
			}
		}
		if err := json.Unmarshal(resp, out); err == nil {
			return nil
		}
		return &APIError{
			HTTPStatus:  status,
			Method:      method,
			RequestPath: requestPath,
			Message:     "invalid response envelope",
			Raw:         resp,
			RequestID:   respHeader.Get("x-request-id"),
		}
	}

	if status >= http.StatusBadRequest || env.Code != "0" {
		return &APIError{
			HTTPStatus:  status,
			Method:      method,
			RequestPath: requestPath,
			Code:        env.Code,
			Message:     env.Msg,
			Raw:         resp,
			RequestID:   respHeader.Get("x-request-id"),
		}
	}

	if out == nil {
		return nil
	}
	if len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	return json.Unmarshal(env.Data, out)
}

// TimeOffset 返回本地时间与服务器时间的偏移。
// 约定：签名时使用 now() - offset。
func (c *Client) TimeOffset() time.Duration {
	return time.Duration(c.timeOffsetNanos.Load())
}
