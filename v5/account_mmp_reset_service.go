package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountMMPResetRequest struct {
	InstType   string `json:"instType,omitempty"`
	InstFamily string `json:"instFamily"`
}

// AccountMMPResetAck 表示重置 MMP 状态返回项。
type AccountMMPResetAck struct {
	Result bool `json:"result"`
}

// AccountMMPResetService 重置 MMP 状态。
type AccountMMPResetService struct {
	c *Client
	r accountMMPResetRequest
}

// NewAccountMMPResetService 创建 AccountMMPResetService。
func (c *Client) NewAccountMMPResetService() *AccountMMPResetService {
	return &AccountMMPResetService{c: c}
}

// InstType 设置交易产品类型（可选，默认 OPTION）。
func (s *AccountMMPResetService) InstType(instType string) *AccountMMPResetService {
	s.r.InstType = instType
	return s
}

// InstFamily 设置交易品种（必填）。
func (s *AccountMMPResetService) InstFamily(instFamily string) *AccountMMPResetService {
	s.r.InstFamily = instFamily
	return s
}

var (
	errAccountMMPResetMissingInstFamily = errors.New("okx: mmp reset requires instFamily")
	errEmptyAccountMMPReset             = errors.New("okx: empty mmp reset response")
)

// Do 重置 MMP 状态（POST /api/v5/account/mmp-reset）。
func (s *AccountMMPResetService) Do(ctx context.Context) (*AccountMMPResetAck, error) {
	if s.r.InstFamily == "" {
		return nil, errAccountMMPResetMissingInstFamily
	}

	var data []AccountMMPResetAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/mmp-reset", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAccountMMPReset
	}
	if !data[0].Result {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/account/mmp-reset",
			RequestID:   requestID,
			Code:        "0",
			Message:     "account mmp reset result is false",
		}
	}
	return &data[0], nil
}
