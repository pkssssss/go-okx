package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// RubikTakerVolumeContractService 获取合约主动买入/卖出情况。
type RubikTakerVolumeContractService struct {
	c *Client

	instId string
	period string
	unit   string
	end    string
	begin  string
	limit  *int
}

// NewRubikTakerVolumeContractService 创建 RubikTakerVolumeContractService。
func (c *Client) NewRubikTakerVolumeContractService() *RubikTakerVolumeContractService {
	return &RubikTakerVolumeContractService{c: c}
}

// InstId 设置产品ID（必填，如 BTC-USDT-SWAP）。
func (s *RubikTakerVolumeContractService) InstId(instId string) *RubikTakerVolumeContractService {
	s.instId = instId
	return s
}

// Period 设置时间粒度（可选）。
func (s *RubikTakerVolumeContractService) Period(period string) *RubikTakerVolumeContractService {
	s.period = period
	return s
}

// Unit 设置买入、卖出的单位（可选）：0=币，1=合约，2=U。
func (s *RubikTakerVolumeContractService) Unit(unit string) *RubikTakerVolumeContractService {
	s.unit = unit
	return s
}

// Begin 设置筛选的开始时间戳（毫秒字符串，可选）。
func (s *RubikTakerVolumeContractService) Begin(begin string) *RubikTakerVolumeContractService {
	s.begin = begin
	return s
}

// End 设置筛选的结束时间戳（毫秒字符串，可选）。
func (s *RubikTakerVolumeContractService) End(end string) *RubikTakerVolumeContractService {
	s.end = end
	return s
}

// Limit 设置分页返回的结果集数量（可选，最大 100）。
func (s *RubikTakerVolumeContractService) Limit(limit int) *RubikTakerVolumeContractService {
	s.limit = &limit
	return s
}

var errRubikTakerVolumeContractMissingInstId = errors.New("okx: rubik taker volume contract requires instId")

// Do 获取合约主动买入/卖出情况（GET /api/v5/rubik/stat/taker-volume-contract）。
func (s *RubikTakerVolumeContractService) Do(ctx context.Context) ([]RubikTakerVolume, error) {
	if s.instId == "" {
		return nil, errRubikTakerVolumeContractMissingInstId
	}

	q := url.Values{}
	q.Set("instId", s.instId)
	if s.period != "" {
		q.Set("period", s.period)
	}
	if s.unit != "" {
		q.Set("unit", s.unit)
	}
	if s.end != "" {
		q.Set("end", s.end)
	}
	if s.begin != "" {
		q.Set("begin", s.begin)
	}
	if s.limit != nil {
		q.Set("limit", strconv.Itoa(*s.limit))
	}

	var data []RubikTakerVolume
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/rubik/stat/taker-volume-contract", q, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
