package okx

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

type accountSetMMPConfigRequest struct {
	InstFamily     string `json:"instFamily"`
	TimeInterval   string `json:"timeInterval"`
	FrozenInterval string `json:"frozenInterval"`
	QtyLimit       string `json:"qtyLimit"`
}

// AccountSetMMPConfigAck 表示设置 MMP 配置返回项。
type AccountSetMMPConfigAck struct {
	FrozenInterval string `json:"frozenInterval"`
	InstFamily     string `json:"instFamily"`
	QtyLimit       string `json:"qtyLimit"`
	TimeInterval   string `json:"timeInterval"`
}

// AccountSetMMPConfigService 设置 MMP 配置。
type AccountSetMMPConfigService struct {
	c *Client
	r accountSetMMPConfigRequest
}

// NewAccountSetMMPConfigService 创建 AccountSetMMPConfigService。
func (c *Client) NewAccountSetMMPConfigService() *AccountSetMMPConfigService {
	return &AccountSetMMPConfigService{c: c}
}

// InstFamily 设置交易品种（必填）。
func (s *AccountSetMMPConfigService) InstFamily(instFamily string) *AccountSetMMPConfigService {
	s.r.InstFamily = instFamily
	return s
}

// TimeInterval 设置时间窗口（毫秒，必填；"0" 代表停用 MMP）。
func (s *AccountSetMMPConfigService) TimeInterval(timeInterval string) *AccountSetMMPConfigService {
	s.r.TimeInterval = timeInterval
	return s
}

// FrozenInterval 设置冻结时间长度（毫秒，必填；"0" 代表一直冻结直到手动重置）。
func (s *AccountSetMMPConfigService) FrozenInterval(frozenInterval string) *AccountSetMMPConfigService {
	s.r.FrozenInterval = frozenInterval
	return s
}

// QtyLimit 设置成交数量上限（必填，需大于 0）。
func (s *AccountSetMMPConfigService) QtyLimit(qtyLimit string) *AccountSetMMPConfigService {
	s.r.QtyLimit = qtyLimit
	return s
}

var (
	errAccountSetMMPConfigMissingRequired = errors.New("okx: set mmp config requires instFamily/timeInterval/frozenInterval/qtyLimit")
	errAccountSetMMPConfigInvalidQtyLimit = errors.New("okx: set mmp config requires qtyLimit > 0")
	errEmptyAccountSetMMPConfig           = errors.New("okx: empty set mmp config response")
	errInvalidAccountSetMMPConfig         = errors.New("okx: invalid set mmp config response")
)

func validateAccountSetMMPConfigAck(ack *AccountSetMMPConfigAck, req accountSetMMPConfigRequest) error {
	if ack == nil || ack.InstFamily == "" || ack.TimeInterval == "" || ack.FrozenInterval == "" || ack.QtyLimit == "" {
		return errInvalidAccountSetMMPConfig
	}
	if ack.InstFamily != req.InstFamily || ack.TimeInterval != req.TimeInterval || ack.FrozenInterval != req.FrozenInterval || ack.QtyLimit != req.QtyLimit {
		return errInvalidAccountSetMMPConfig
	}
	return nil
}

// Do 设置 MMP 配置（POST /api/v5/account/mmp-config）。
func (s *AccountSetMMPConfigService) Do(ctx context.Context) (*AccountSetMMPConfigAck, error) {
	if s.r.InstFamily == "" || s.r.TimeInterval == "" || s.r.FrozenInterval == "" || s.r.QtyLimit == "" {
		return nil, errAccountSetMMPConfigMissingRequired
	}

	qty, err := strconv.ParseInt(s.r.QtyLimit, 10, 64)
	if err != nil || qty <= 0 {
		return nil, errAccountSetMMPConfigInvalidQtyLimit
	}

	var data []AccountSetMMPConfigAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/mmp-config", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/mmp-config", requestID, errEmptyAccountSetMMPConfig)
	}
	if err := validateAccountSetMMPConfigAck(&data[0], s.r); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/mmp-config", requestID, err)
	}
	return &data[0], nil
}
