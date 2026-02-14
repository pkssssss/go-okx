package okx

import (
	"context"
	"errors"
	"net/http"
)

// AccountActivateOptionAck 表示开通期权交易返回项。
type AccountActivateOptionAck struct {
	TS UnixMilli `json:"ts"`
}

// AccountActivateOptionService 开通期权交易。
type AccountActivateOptionService struct {
	c *Client
}

// NewAccountActivateOptionService 创建 AccountActivateOptionService。
func (c *Client) NewAccountActivateOptionService() *AccountActivateOptionService {
	return &AccountActivateOptionService{c: c}
}

var errEmptyAccountActivateOption = errors.New("okx: empty activate option response")

// Do 开通期权交易（POST /api/v5/account/activate-option）。
func (s *AccountActivateOptionService) Do(ctx context.Context) (*AccountActivateOptionAck, error) {
	var data []AccountActivateOptionAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/activate-option", nil, nil, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/activate-option", requestID, errEmptyAccountActivateOption)
	}
	return &data[0], nil
}
