package okx

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

// SprdCancelAllAfterAck 表示设置倒计时撤单的返回项（价差交易）。
type SprdCancelAllAfterAck struct {
	TriggerTime int64 `json:"triggerTime,string"`
	TS          int64 `json:"ts,string"`
}

// SprdCancelAllAfterService 设置倒计时全部撤单（价差交易）。
type SprdCancelAllAfterService struct {
	c *Client

	timeOut string
}

// NewSprdCancelAllAfterService 创建 SprdCancelAllAfterService。
func (c *Client) NewSprdCancelAllAfterService() *SprdCancelAllAfterService {
	return &SprdCancelAllAfterService{c: c}
}

// TimeOut 设置倒计时秒数（必填：0 或 10-120）。
func (s *SprdCancelAllAfterService) TimeOut(timeOut string) *SprdCancelAllAfterService {
	s.timeOut = timeOut
	return s
}

var (
	errSprdCancelAllAfterMissingTimeOut = errors.New("okx: sprd cancel all after requires timeOut")
	errSprdCancelAllAfterInvalidTimeOut = errors.New("okx: sprd cancel all after requires timeOut=0 or 10-120")
	errEmptySprdCancelAllAfterResponse  = errors.New("okx: empty sprd cancel all after response")
)

type sprdCancelAllAfterRequest struct {
	TimeOut string `json:"timeOut"`
}

// Do 设置倒计时全部撤单（POST /api/v5/sprd/cancel-all-after）。
func (s *SprdCancelAllAfterService) Do(ctx context.Context) (*SprdCancelAllAfterAck, error) {
	if s.timeOut == "" {
		return nil, errSprdCancelAllAfterMissingTimeOut
	}

	timeoutSeconds, err := strconv.Atoi(s.timeOut)
	if err != nil || (timeoutSeconds != 0 && (timeoutSeconds < 10 || timeoutSeconds > 120)) {
		return nil, errSprdCancelAllAfterInvalidTimeOut
	}

	req := sprdCancelAllAfterRequest{TimeOut: s.timeOut}

	var data []SprdCancelAllAfterAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/sprd/cancel-all-after", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySprdCancelAllAfterResponse
	}
	return &data[0], nil
}
