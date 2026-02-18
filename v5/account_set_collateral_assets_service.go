package okx

import (
	"context"
	"errors"
	"net/http"
)

type accountSetCollateralAssetsRequest struct {
	Type string `json:"type"`

	CcyList []string `json:"ccyList,omitempty"`

	CollateralEnabled *bool `json:"collateralEnabled"`
}

// AccountSetCollateralAssetsAck 表示设置质押币种返回项。
type AccountSetCollateralAssetsAck struct {
	Type              string   `json:"type"`
	CcyList           []string `json:"ccyList"`
	CollateralEnabled bool     `json:"collateralEnabled"`
}

// AccountSetCollateralAssetsService 设置质押币种。
type AccountSetCollateralAssetsService struct {
	c *Client
	r accountSetCollateralAssetsRequest
}

// NewAccountSetCollateralAssetsService 创建 AccountSetCollateralAssetsService。
func (c *Client) NewAccountSetCollateralAssetsService() *AccountSetCollateralAssetsService {
	return &AccountSetCollateralAssetsService{c: c}
}

// Type 设置币种类型（必填）：all/custom。
func (s *AccountSetCollateralAssetsService) Type(typ string) *AccountSetCollateralAssetsService {
	s.r.Type = typ
	return s
}

// CollateralEnabled 设置是否为质押币种（必填）。
func (s *AccountSetCollateralAssetsService) CollateralEnabled(enable bool) *AccountSetCollateralAssetsService {
	s.r.CollateralEnabled = &enable
	return s
}

// CcyList 设置币种列表（type=custom 时必填）。
func (s *AccountSetCollateralAssetsService) CcyList(ccys []string) *AccountSetCollateralAssetsService {
	s.r.CcyList = ccys
	return s
}

var (
	errAccountSetCollateralAssetsMissingRequired = errors.New("okx: set collateral assets requires type/collateralEnabled")
	errAccountSetCollateralAssetsMissingCcyList  = errors.New("okx: set collateral assets with type=custom requires ccyList")
	errEmptyAccountSetCollateralAssets           = errors.New("okx: empty set collateral assets response")
	errInvalidAccountSetCollateralAssets         = errors.New("okx: invalid set collateral assets response")
)

func validateAccountSetCollateralAssetsAck(ack *AccountSetCollateralAssetsAck, req accountSetCollateralAssetsRequest) error {
	if ack == nil || ack.Type == "" {
		return errInvalidAccountSetCollateralAssets
	}
	if ack.Type != req.Type {
		return errInvalidAccountSetCollateralAssets
	}
	if req.Type == "custom" && len(ack.CcyList) == 0 {
		return errInvalidAccountSetCollateralAssets
	}
	return nil
}

// Do 设置质押币种（POST /api/v5/account/set-collateral-assets）。
func (s *AccountSetCollateralAssetsService) Do(ctx context.Context) (*AccountSetCollateralAssetsAck, error) {
	if s.r.Type == "" || s.r.CollateralEnabled == nil {
		return nil, errAccountSetCollateralAssetsMissingRequired
	}
	if s.r.Type == "custom" && len(s.r.CcyList) == 0 {
		return nil, errAccountSetCollateralAssetsMissingCcyList
	}

	var data []AccountSetCollateralAssetsAck
	requestID, err := s.c.doWithHeadersAndRequestID(ctx, http.MethodPost, "/api/v5/account/set-collateral-assets", nil, s.r, true, nil, &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, newEmptyDataAPIError(http.MethodPost, "/api/v5/account/set-collateral-assets", requestID, errEmptyAccountSetCollateralAssets)
	}
	if err := validateAccountSetCollateralAssetsAck(&data[0], s.r); err != nil {
		return nil, newInvalidDataAPIError(http.MethodPost, "/api/v5/account/set-collateral-assets", requestID, err)
	}
	return &data[0], nil
}
