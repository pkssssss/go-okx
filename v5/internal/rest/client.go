package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// Client 只负责最基础的 HTTP 发送与响应读取，不包含 OKX 业务语义（签名/解包/错误码等）。
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

func (c *Client) Do(ctx context.Context, method, requestPath string, body []byte, header http.Header) (status int, resp []byte, respHeader http.Header, err error) {
	fullURL := c.BaseURL + requestPath

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
