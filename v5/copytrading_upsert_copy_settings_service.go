package okx

import (
	"context"
	"errors"
	"net/http"
)

// CopyTradingUpsertCopySettingsService 首次跟单设置/修改跟单设置。
//
// 说明：该 Service 通过不同构造函数区分 endpoint：
// - NewCopyTradingFirstCopySettingsService => POST /api/v5/copytrading/first-copy-settings
// - NewCopyTradingAmendCopySettingsService => POST /api/v5/copytrading/amend-copy-settings
type CopyTradingUpsertCopySettingsService struct {
	c        *Client
	endpoint string

	instType        string
	uniqueCode      string
	copyMgnMode     string
	copyInstIdType  string
	instId          string
	copyMode        string
	copyTotalAmt    string
	copyAmt         string
	copyRatio       string
	tpRatio         string
	slRatio         string
	slTotalAmt      string
	subPosCloseType string
	tag             string
}

// NewCopyTradingFirstCopySettingsService 创建首次跟单设置 Service。
func (c *Client) NewCopyTradingFirstCopySettingsService() *CopyTradingUpsertCopySettingsService {
	return &CopyTradingUpsertCopySettingsService{c: c, endpoint: "/api/v5/copytrading/first-copy-settings"}
}

// NewCopyTradingAmendCopySettingsService 创建修改跟单设置 Service。
func (c *Client) NewCopyTradingAmendCopySettingsService() *CopyTradingUpsertCopySettingsService {
	return &CopyTradingUpsertCopySettingsService{c: c, endpoint: "/api/v5/copytrading/amend-copy-settings"}
}

// InstType 设置产品类型（默认 SWAP）。
func (s *CopyTradingUpsertCopySettingsService) InstType(instType string) *CopyTradingUpsertCopySettingsService {
	s.instType = instType
	return s
}

// UniqueCode 设置交易员唯一标识码（必填）。
func (s *CopyTradingUpsertCopySettingsService) UniqueCode(uniqueCode string) *CopyTradingUpsertCopySettingsService {
	s.uniqueCode = uniqueCode
	return s
}

// CopyMgnMode 设置跟单保证金模式（cross/isolated/copy；必填）。
func (s *CopyTradingUpsertCopySettingsService) CopyMgnMode(copyMgnMode string) *CopyTradingUpsertCopySettingsService {
	s.copyMgnMode = copyMgnMode
	return s
}

// CopyInstIdType 设置跟单合约设置类型（custom/copy；必填；custom 时 instId 必填）。
func (s *CopyTradingUpsertCopySettingsService) CopyInstIdType(copyInstIdType string) *CopyTradingUpsertCopySettingsService {
	s.copyInstIdType = copyInstIdType
	return s
}

// InstId 设置产品 ID（copyInstIdType=custom 时必填；可传入多条，以逗号区分）。
func (s *CopyTradingUpsertCopySettingsService) InstId(instId string) *CopyTradingUpsertCopySettingsService {
	s.instId = instId
	return s
}

// CopyMode 设置跟单模式（fixed_amount/ratio_copy；默认 fixed_amount）。
func (s *CopyTradingUpsertCopySettingsService) CopyMode(copyMode string) *CopyTradingUpsertCopySettingsService {
	s.copyMode = copyMode
	return s
}

// CopyTotalAmt 设置投入的最大跟单金额（USDT；必填）。
func (s *CopyTradingUpsertCopySettingsService) CopyTotalAmt(copyTotalAmt string) *CopyTradingUpsertCopySettingsService {
	s.copyTotalAmt = copyTotalAmt
	return s
}

// CopyAmt 设置单笔跟随金额（USDT；fixed_amount 模式必填）。
func (s *CopyTradingUpsertCopySettingsService) CopyAmt(copyAmt string) *CopyTradingUpsertCopySettingsService {
	s.copyAmt = copyAmt
	return s
}

// CopyRatio 设置跟单比例（ratio_copy 模式必填）。
func (s *CopyTradingUpsertCopySettingsService) CopyRatio(copyRatio string) *CopyTradingUpsertCopySettingsService {
	s.copyRatio = copyRatio
	return s
}

// TpRatio 设置单笔止盈百分比（0.1 代表 10%）。
func (s *CopyTradingUpsertCopySettingsService) TpRatio(tpRatio string) *CopyTradingUpsertCopySettingsService {
	s.tpRatio = tpRatio
	return s
}

// SlRatio 设置单笔止损百分比（0.1 代表 10%）。
func (s *CopyTradingUpsertCopySettingsService) SlRatio(slRatio string) *CopyTradingUpsertCopySettingsService {
	s.slRatio = slRatio
	return s
}

// SlTotalAmt 设置跟单止损总金额（USDT）。
func (s *CopyTradingUpsertCopySettingsService) SlTotalAmt(slTotalAmt string) *CopyTradingUpsertCopySettingsService {
	s.slTotalAmt = slTotalAmt
	return s
}

// SubPosCloseType 设置剩余仓位处理方式（market_close/copy_close/manual_close；必填）。
func (s *CopyTradingUpsertCopySettingsService) SubPosCloseType(subPosCloseType string) *CopyTradingUpsertCopySettingsService {
	s.subPosCloseType = subPosCloseType
	return s
}

// Tag 设置订单标签（1-16）。
func (s *CopyTradingUpsertCopySettingsService) Tag(tag string) *CopyTradingUpsertCopySettingsService {
	s.tag = tag
	return s
}

var (
	errCopyTradingUpsertCopySettingsInvalidEndpoint  = errors.New("okx: invalid copytrading copy settings endpoint")
	errCopyTradingUpsertCopySettingsMissingRequired  = errors.New("okx: copytrading copy settings requires uniqueCode/copyMgnMode/copyInstIdType/copyTotalAmt/subPosCloseType")
	errCopyTradingUpsertCopySettingsMissingInstId    = errors.New("okx: copytrading copy settings requires instId for copyInstIdType=custom")
	errCopyTradingUpsertCopySettingsMissingCopyAmt   = errors.New("okx: copytrading copy settings requires copyAmt for fixed_amount")
	errCopyTradingUpsertCopySettingsMissingCopyRatio = errors.New("okx: copytrading copy settings requires copyRatio for ratio_copy")
	errEmptyCopyTradingUpsertCopySettingsResponse    = errors.New("okx: empty copytrading copy settings response")
)

type copyTradingUpsertCopySettingsRequest struct {
	InstType        string `json:"instType,omitempty"`
	UniqueCode      string `json:"uniqueCode"`
	CopyMgnMode     string `json:"copyMgnMode"`
	CopyInstIdType  string `json:"copyInstIdType"`
	InstId          string `json:"instId,omitempty"`
	CopyMode        string `json:"copyMode,omitempty"`
	CopyTotalAmt    string `json:"copyTotalAmt"`
	CopyAmt         string `json:"copyAmt,omitempty"`
	CopyRatio       string `json:"copyRatio,omitempty"`
	TpRatio         string `json:"tpRatio,omitempty"`
	SlRatio         string `json:"slRatio,omitempty"`
	SlTotalAmt      string `json:"slTotalAmt,omitempty"`
	SubPosCloseType string `json:"subPosCloseType"`
	Tag             string `json:"tag,omitempty"`
}

// Do 首次/修改跟单设置（POST /api/v5/copytrading/first-copy-settings 或 /api/v5/copytrading/amend-copy-settings）。
func (s *CopyTradingUpsertCopySettingsService) Do(ctx context.Context) (*CopyTradingResult, error) {
	if s.endpoint == "" {
		return nil, errCopyTradingUpsertCopySettingsInvalidEndpoint
	}
	if s.uniqueCode == "" || s.copyMgnMode == "" || s.copyInstIdType == "" || s.copyTotalAmt == "" || s.subPosCloseType == "" {
		return nil, errCopyTradingUpsertCopySettingsMissingRequired
	}
	if s.copyInstIdType == "custom" && s.instId == "" {
		return nil, errCopyTradingUpsertCopySettingsMissingInstId
	}

	copyMode := s.copyMode
	if copyMode == "" {
		copyMode = "fixed_amount"
	}
	switch copyMode {
	case "fixed_amount":
		if s.copyAmt == "" {
			return nil, errCopyTradingUpsertCopySettingsMissingCopyAmt
		}
	case "ratio_copy":
		if s.copyRatio == "" {
			return nil, errCopyTradingUpsertCopySettingsMissingCopyRatio
		}
	}

	req := copyTradingUpsertCopySettingsRequest{
		InstType:        s.instType,
		UniqueCode:      s.uniqueCode,
		CopyMgnMode:     s.copyMgnMode,
		CopyInstIdType:  s.copyInstIdType,
		InstId:          s.instId,
		CopyMode:        s.copyMode,
		CopyTotalAmt:    s.copyTotalAmt,
		CopyAmt:         s.copyAmt,
		CopyRatio:       s.copyRatio,
		TpRatio:         s.tpRatio,
		SlRatio:         s.slRatio,
		SlTotalAmt:      s.slTotalAmt,
		SubPosCloseType: s.subPosCloseType,
		Tag:             s.tag,
	}

	var data []CopyTradingResult
	if err := s.c.do(ctx, http.MethodPost, s.endpoint, nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyCopyTradingUpsertCopySettingsResponse
	}
	return &data[0], nil
}
