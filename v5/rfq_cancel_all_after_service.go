package okx

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

// RFQCancelAllAfterAck 表示设置倒计时取消所有报价单的返回项。
type RFQCancelAllAfterAck struct {
	TriggerTime int64 `json:"triggerTime,string"`
	TS          int64 `json:"ts,string"`
}

// RFQCancelAllAfterService 设置倒计时取消所有报价单。
type RFQCancelAllAfterService struct {
	c *Client

	timeOut string
}

// NewRFQCancelAllAfterService 创建 RFQCancelAllAfterService。
func (c *Client) NewRFQCancelAllAfterService() *RFQCancelAllAfterService {
	return &RFQCancelAllAfterService{c: c}
}

// TimeOut 设置倒计时秒数（必填：0 或 10-120）。
func (s *RFQCancelAllAfterService) TimeOut(timeOut string) *RFQCancelAllAfterService {
	s.timeOut = timeOut
	return s
}

var (
	errRFQCancelAllAfterMissingTimeOut = errors.New("okx: rfq cancel all after requires timeOut")
	errRFQCancelAllAfterInvalidTimeOut = errors.New("okx: rfq cancel all after requires timeOut=0 or 10-120")
	errEmptyRFQCancelAllAfterResponse  = errors.New("okx: empty rfq cancel all after response")
)

type rfqCancelAllAfterRequest struct {
	TimeOut string `json:"timeOut"`
}

// Do 设置倒计时取消所有报价单（POST /api/v5/rfq/cancel-all-after）。
func (s *RFQCancelAllAfterService) Do(ctx context.Context) (*RFQCancelAllAfterAck, error) {
	if s.timeOut == "" {
		return nil, errRFQCancelAllAfterMissingTimeOut
	}

	timeoutSeconds, err := strconv.Atoi(s.timeOut)
	if err != nil || (timeoutSeconds != 0 && (timeoutSeconds < 10 || timeoutSeconds > 120)) {
		return nil, errRFQCancelAllAfterInvalidTimeOut
	}

	req := rfqCancelAllAfterRequest{TimeOut: s.timeOut}

	var data []RFQCancelAllAfterAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/rfq/cancel-all-after", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/rfq/cancel-all-after", requestID, errEmptyRFQCancelAllAfterResponse)
	}
	return &data[0], nil
}
