package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingSetInstrumentsService 交易员修改带单产品。
type CopyTradingSetInstrumentsService struct {
	c *Client

	instType string
	instId   string
}

// NewCopyTradingSetInstrumentsService 创建 CopyTradingSetInstrumentsService。
func (c *Client) NewCopyTradingSetInstrumentsService() *CopyTradingSetInstrumentsService {
	return &CopyTradingSetInstrumentsService{c: c}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingSetInstrumentsService) InstType(instType string) *CopyTradingSetInstrumentsService {
	s.instType = instType
	return s
}

// InstId 设置产品 ID（必填；多个用半角逗号分隔）。
func (s *CopyTradingSetInstrumentsService) InstId(instId string) *CopyTradingSetInstrumentsService {
	s.instId = instId
	return s
}

var (
	errCopyTradingSetInstrumentsMissingInstId = errors.New("okx: copytrading set instruments requires instId")
	errEmptyCopyTradingSetInstrumentsResponse = errors.New("okx: empty copytrading set instruments response")
)

type copyTradingSetInstrumentsRequest struct {
	InstType string `json:"instType,omitempty"`
	InstId   string `json:"instId"`
}

// Do 修改带单产品（POST /api/v5/copytrading/set-instruments）。
func (s *CopyTradingSetInstrumentsService) Do(ctx context.Context) ([]CopyTradingInstrument, error) {
	if s.instId == "" {
		return nil, errCopyTradingSetInstrumentsMissingInstId
	}

	req := copyTradingSetInstrumentsRequest{
		InstType: s.instType,
		InstId:   s.instId,
	}

	var data []CopyTradingInstrument
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/copytrading/set-instruments", nil, req, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/copytrading/set-instruments", requestID, errEmptyCopyTradingSetInstrumentsResponse)
	}
	return data, nil
}
