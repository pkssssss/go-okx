package okx

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

// TradeMassCancelAck 表示撤销 MMP 订单返回项。
type TradeMassCancelAck struct {
	Result bool `json:"result"`
}

// MassCancelService 撤销同一交易品种下用户所有的 MMP 挂单。
type MassCancelService struct {
	c *Client

	instType     string
	instFamily   string
	lockInterval string
}

// NewMassCancelService 创建 MassCancelService。
func (c *Client) NewMassCancelService() *MassCancelService {
	return &MassCancelService{c: c}
}

// InstType 设置交易产品类型（必填；目前仅支持 OPTION）。
func (s *MassCancelService) InstType(instType string) *MassCancelService {
	s.instType = instType
	return s
}

// InstFamily 设置交易品种（必填，如 BTC-USD）。
func (s *MassCancelService) InstFamily(instFamily string) *MassCancelService {
	s.instFamily = instFamily
	return s
}

// LockInterval 设置锁定时长（毫秒，可选；范围 0-10000）。
func (s *MassCancelService) LockInterval(lockInterval string) *MassCancelService {
	s.lockInterval = lockInterval
	return s
}

var (
	errMassCancelMissingRequired   = errors.New("okx: mass cancel requires instType/instFamily")
	errMassCancelInvalidLockWindow = errors.New("okx: mass cancel requires lockInterval in [0, 10000]")
	errEmptyMassCancelResponse     = errors.New("okx: empty mass cancel response")
)

type massCancelRequest struct {
	InstType     string `json:"instType"`
	InstFamily   string `json:"instFamily"`
	LockInterval string `json:"lockInterval,omitempty"`
}

// Do 撤销 MMP 订单（POST /api/v5/trade/mass-cancel）。
func (s *MassCancelService) Do(ctx context.Context) (*TradeMassCancelAck, error) {
	if s.instType == "" || s.instFamily == "" {
		return nil, errMassCancelMissingRequired
	}
	if s.lockInterval != "" {
		n, err := strconv.Atoi(s.lockInterval)
		if err != nil || n < 0 || n > 10_000 {
			return nil, errMassCancelInvalidLockWindow
		}
	}

	req := massCancelRequest{
		InstType:     s.instType,
		InstFamily:   s.instFamily,
		LockInterval: s.lockInterval,
	}

	var data []TradeMassCancelAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/mass-cancel", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/mass-cancel", requestID, errEmptyMassCancelResponse)
	}
	if !data[0].Result {
		return nil, &APIError{
			HTTPStatus:  http.StatusOK,
			Method:      http.MethodPost,
			RequestPath: "/api/v5/trade/mass-cancel",
			RequestID:   requestID,
			Code:        "0",
			Message:     "mass cancel result is false",
		}
	}
	return &data[0], nil
}
