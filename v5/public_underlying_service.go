package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// PublicUnderlyingService 获取衍生品标的指数。
type PublicUnderlyingService struct {
	c *Client

	instType string
}

// NewPublicUnderlyingService 创建 PublicUnderlyingService。
func (c *Client) NewPublicUnderlyingService() *PublicUnderlyingService {
	return &PublicUnderlyingService{c: c}
}

// InstType 设置产品类型（必填：SWAP/FUTURES/OPTION）。
func (s *PublicUnderlyingService) InstType(instType string) *PublicUnderlyingService {
	s.instType = instType
	return s
}

var errPublicUnderlyingMissingInstType = errors.New("okx: public underlying requires instType")

// Do 获取衍生品标的指数（GET /api/v5/public/underlying）。
func (s *PublicUnderlyingService) Do(ctx context.Context) ([]string, error) {
	if s.instType == "" {
		return nil, errPublicUnderlyingMissingInstType
	}

	q := url.Values{}
	q.Set("instType", s.instType)

	var raw [][]string
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/public/underlying", q, nil, false, &raw); err != nil {
		return nil, err
	}

	var out []string
	for _, group := range raw {
		out = append(out, group...)
	}
	return out, nil
}
