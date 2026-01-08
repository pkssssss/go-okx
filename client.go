package okx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/pkssssss/go-okx/internal/rest"
	"github.com/pkssssss/go-okx/internal/sign"
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

// Client 是 OKX V5 API 客户端（REST + WS）。
// v0.1 先落地 REST 核心能力，WS 随后补齐。
type Client struct {
	rest *rest.Client

	creds *Credentials
	demo  bool

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
		now: time.Now,
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

	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	if c.demo {
		header.Set("x-simulated-trading", "1")
	}

	if signed {
		if c.creds == nil || c.creds.APIKey == "" || c.creds.SecretKey == "" || c.creds.Passphrase == "" {
			return errMissingCredentials
		}

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

	status, resp, respHeader, err := c.rest.Do(ctx, method, requestPath, bodyBytes, header)
	if err != nil {
		return err
	}

	var env responseEnvelope
	if err := json.Unmarshal(resp, &env); err != nil {
		return &APIError{
			HTTPStatus:  status,
			Method:      method,
			RequestPath: requestPath,
			Message:     "invalid JSON response",
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
