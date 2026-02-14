package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client 只负责最基础的 HTTP 发送与响应读取，不包含 OKX 业务语义（签名/解包/错误码等）。
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	// DefaultTimeout 用于在 ctx 未设置 deadline 时，给请求提供一个安全的默认超时（Fail-Fast）。
	// 若为 0，则使用内部默认值；若为负数，则表示禁用默认超时（不建议）。
	DefaultTimeout time.Duration
	// MaxResponseBodyBytes 用于限制响应体最大读取字节数，避免异常响应导致内存/延迟尖刺。
	// 若为 0，则使用内部默认值；若为负数，则表示禁用上限（不建议）。
	MaxResponseBodyBytes int64
}

// ResponseBodyTooLargeError 表示响应体超过预期上限（通常是上游异常/代理干扰/错误返回）。
type ResponseBodyTooLargeError struct {
	Method      string
	RequestPath string
	MaxBytes    int64
}

func (e *ResponseBodyTooLargeError) Error() string {
	if e == nil {
		return "rest: response body too large"
	}
	return "rest: response body too large"
}

var (
	errRedirectBlockedSignedRequest = errors.New("rest: redirect blocked for signed request")
	errRedirectBlockedCrossHost     = errors.New("rest: redirect to different host blocked")
)

var defaultHTTPClient = &http.Client{
	CheckRedirect: checkRedirect,
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) == 0 {
		return nil
	}
	from := via[0]
	if from != nil && hasOKAccessHeaders(from.Header) {
		return errRedirectBlockedSignedRequest
	}
	if from != nil && from.URL != nil && req != nil && req.URL != nil {
		if !sameSchemeAndHost(from.URL, req.URL) {
			return errRedirectBlockedCrossHost
		}
	}
	return nil
}

func composeCheckRedirect(userCheckRedirect func(req *http.Request, via []*http.Request) error) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if err := checkRedirect(req, via); err != nil {
			return err
		}
		if userCheckRedirect != nil {
			return userCheckRedirect(req, via)
		}
		return nil
	}
}

func hasOKAccessHeaders(h http.Header) bool {
	if h == nil {
		return false
	}
	return h.Get("OK-ACCESS-KEY") != "" ||
		h.Get("OK-ACCESS-PASSPHRASE") != "" ||
		h.Get("OK-ACCESS-TIMESTAMP") != "" ||
		h.Get("OK-ACCESS-SIGN") != ""
}

func sameSchemeAndHost(a, b *url.URL) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Scheme == b.Scheme && a.Host == b.Host
}

// ContextWithDefaultTimeout 在 ctx 未设置 deadline 时，为其附加 DefaultTimeout（Fail-Fast）。
// 返回的 cancel 需要由调用方负责调用（若为 nil 则无需调用）。
func (c *Client) ContextWithDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, nil
	}

	timeout := c.DefaultTimeout
	switch {
	case timeout == 0:
		timeout = 10 * time.Second
	case timeout < 0:
		timeout = 0
	}
	if timeout <= 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, timeout)
}

func (c *Client) Do(ctx context.Context, method, requestPath string, body []byte, header http.Header) (status int, resp []byte, respHeader http.Header, err error) {
	fullURL := c.BaseURL + requestPath

	ctx, cancel := c.ContextWithDefaultTimeout(ctx)
	if cancel != nil {
		defer cancel()
	}

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return 0, nil, nil, err
	}

	if header != nil {
		req.Header = header.Clone()
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	hc := c.HTTPClient
	switch {
	case hc == nil:
		hc = defaultHTTPClient
	default:
		hcCopy := *hc
		hcCopy.CheckRedirect = composeCheckRedirect(hc.CheckRedirect)
		hc = &hcCopy
	}

	res, err := hc.Do(req)
	if err != nil {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
		return 0, nil, nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	maxBody := c.MaxResponseBodyBytes
	switch {
	case maxBody == 0:
		maxBody = 16 << 20 // 16MiB：覆盖 OKX 常规响应并限制异常大 body 风险
	case maxBody < 0:
		maxBody = 0
	}

	reader := io.Reader(res.Body)
	if maxBody > 0 {
		reader = io.LimitReader(res.Body, maxBody+1)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return 0, nil, nil, err
	}
	if maxBody > 0 && int64(len(data)) > maxBody {
		return 0, nil, nil, &ResponseBodyTooLargeError{Method: method, RequestPath: requestPath, MaxBytes: maxBody}
	}

	return res.StatusCode, data, res.Header.Clone(), nil
}

// BuildRequestPath 把 endpoint 与 query 编码为 OKX 使用的 requestPath（用于签名与实际请求）。
func BuildRequestPath(endpoint string, query url.Values) string {
	if query == nil {
		return endpoint
	}
	qs := query.Encode()
	if qs == "" {
		return endpoint
	}
	return endpoint + "?" + qs
}
