package rest

import (
	"bytes"
	"context"
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
}

func (c *Client) Do(ctx context.Context, method, requestPath string, body []byte, header http.Header) (status int, resp []byte, respHeader http.Header, err error) {
	fullURL := c.BaseURL + requestPath

	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Deadline(); !ok {
		timeout := c.DefaultTimeout
		switch {
		case timeout == 0:
			timeout = 10 * time.Second
		case timeout < 0:
			timeout = 0
		}
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
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
	if hc == nil {
		hc = http.DefaultClient
	}

	res, err := hc.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, nil, err
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
