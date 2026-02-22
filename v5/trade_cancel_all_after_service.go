package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// TradeCancelAllAfterAck 表示设置倒计时撤单的返回项。
type TradeCancelAllAfterAck struct {
	TriggerTime int64  `json:"triggerTime,string"`
	Tag         string `json:"tag"`
	TS          int64  `json:"ts,string"`
}

// CancelAllAfterService 设置倒计时全部撤单。
type CancelAllAfterService struct {
	c *Client

	timeOut string
	tag     string
}

// NewCancelAllAfterService 创建 CancelAllAfterService。
func (c *Client) NewCancelAllAfterService() *CancelAllAfterService {
	return &CancelAllAfterService{c: c}
}

// TimeOut 设置倒计时秒数（必填：0 或 10-120）。
func (s *CancelAllAfterService) TimeOut(timeOut string) *CancelAllAfterService {
	s.timeOut = timeOut
	return s
}

// Tag 设置订单来源（可选，用于标识）。
func (s *CancelAllAfterService) Tag(tag string) *CancelAllAfterService {
	s.tag = tag
	return s
}

var (
	errCancelAllAfterMissingTimeOut  = errors.New("okx: cancel all after requires timeOut")
	errCancelAllAfterInvalidTimeOut  = errors.New("okx: cancel all after requires timeOut=0 or 10-120")
	errEmptyCancelAllAfterResponse   = errors.New("okx: empty cancel all after response")
	errInvalidCancelAllAfterResponse = errors.New("okx: invalid cancel all after response")
)

type cancelAllAfterRequest struct {
	TimeOut string `json:"timeOut"`
	Tag     string `json:"tag,omitempty"`
}

func validateCancelAllAfterAck(ack *TradeCancelAllAfterAck, req cancelAllAfterRequest) error {
	if ack == nil || ack.TriggerTime <= 0 || ack.TS <= 0 {
		return errInvalidCancelAllAfterResponse
	}
	if req.Tag != "" && ack.Tag != req.Tag {
		return errInvalidCancelAllAfterResponse
	}
	return nil
}

// Do 设置倒计时全部撤单（POST /api/v5/trade/cancel-all-after）。
func (s *CancelAllAfterService) Do(ctx context.Context) (*TradeCancelAllAfterAck, error) {
	if s.timeOut == "" {
		return nil, errCancelAllAfterMissingTimeOut
	}

	timeoutSeconds, err := strconv.Atoi(s.timeOut)
	if err != nil || (timeoutSeconds != 0 && (timeoutSeconds < 10 || timeoutSeconds > 120)) {
		return nil, errCancelAllAfterInvalidTimeOut
	}

	req := cancelAllAfterRequest{
		TimeOut: s.timeOut,
		Tag:     s.tag,
	}

	var data []TradeCancelAllAfterAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/trade/cancel-all-after", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/trade/cancel-all-after", requestID, errEmptyCancelAllAfterResponse)
	}
	if len(data) != 1 {
		return nil, newInvalidDataAPIError(
			http.MethodPost,
			"/api/v5/trade/cancel-all-after",
			requestID,
			fmt.Errorf("%w: expected 1 ack, got %d", errInvalidCancelAllAfterResponse, len(data)),
		)
	}
	if err := validateCancelAllAfterAck(&data[0], req); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/trade/cancel-all-after", requestID, err)
	}
	return &data[0], nil
}
