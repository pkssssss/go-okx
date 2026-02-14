package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetGreeksRequest struct {
	GreeksType string `json:"greeksType"`
}

// AccountSetGreeksAck 表示期权 greeks 展示方式切换返回项。
type AccountSetGreeksAck struct {
	GreeksType string `json:"greeksType"`
}

// AccountSetGreeksService 期权 greeks 的 PA/BS 切换。
type AccountSetGreeksService struct {
	c   *Client
	req accountSetGreeksRequest
}

// NewAccountSetGreeksService 创建 AccountSetGreeksService。
func (c *Client) NewAccountSetGreeksService() *AccountSetGreeksService {
	return &AccountSetGreeksService{c: c}
}

// GreeksType 设置希腊字母展示方式（必填：PA/BS）。
func (s *AccountSetGreeksService) GreeksType(greeksType string) *AccountSetGreeksService {
	s.req.GreeksType = greeksType
	return s
}

var (
	errAccountSetGreeksMissingGreeksType = errors.New("okx: set greeks requires greeksType")
	errEmptyAccountSetGreeks             = errors.New("okx: empty set greeks response")
	errInvalidAccountSetGreeks           = errors.New("okx: invalid set greeks response")
)

func validateAccountSetGreeksAck(ack *AccountSetGreeksAck, req accountSetGreeksRequest) error {
	if ack == nil || ack.GreeksType == "" || ack.GreeksType != req.GreeksType {
		return errInvalidAccountSetGreeks
	}
	return nil
}

// Do 切换期权 greeks 展示方式（POST /api/v5/account/set-greeks）。
func (s *AccountSetGreeksService) Do(ctx context.Context) (*AccountSetGreeksAck, error) {
	if s.req.GreeksType == "" {
		return nil, errAccountSetGreeksMissingGreeksType
	}

	var data []AccountSetGreeksAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-greeks", nil, s.req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-greeks", requestID, errEmptyAccountSetGreeks)
	}
	if err := validateAccountSetGreeksAck(&data[0], s.req); err != nil {
		return nil, err
	}
	return &data[0], nil
}
