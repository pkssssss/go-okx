package okx

import (
	"context"
	"errors"
	"net/http"
)

// EasyConvertService 进行小币一键兑换主流币交易。
type EasyConvertService struct {
	c *Client

	fromCcy []string
	toCcy   string
	source  string
}

// NewEasyConvertService 创建 EasyConvertService。
func (c *Client) NewEasyConvertService() *EasyConvertService {
	return &EasyConvertService{c: c}
}

// FromCcy 设置小币支付币种列表（必填；单次最多 5 个）。
func (s *EasyConvertService) FromCcy(fromCcy []string) *EasyConvertService {
	s.fromCcy = fromCcy
	return s
}

// ToCcy 设置兑换的主流币（必填；不能与 fromCcy 重复）。
func (s *EasyConvertService) ToCcy(toCcy string) *EasyConvertService {
	s.toCcy = toCcy
	return s
}

// Source 设置资金来源（可选：1=交易账户，2=资金账户）。
func (s *EasyConvertService) Source(source string) *EasyConvertService {
	s.source = source
	return s
}

var (
	errEasyConvertMissingRequired = errors.New("okx: easy convert requires fromCcy/toCcy")
	errEasyConvertTooManyFromCcy  = errors.New("okx: easy convert supports up to 5 fromCcy")
	errEasyConvertEmptyFromCcy    = errors.New("okx: easy convert requires non-empty fromCcy items")
	errEasyConvertSameCurrency    = errors.New("okx: easy convert requires toCcy not in fromCcy")
	errEmptyEasyConvertResponse   = errors.New("okx: empty easy convert response")
)

type easyConvertRequest struct {
	FromCcy []string `json:"fromCcy"`
	ToCcy   string   `json:"toCcy"`
	Source  string   `json:"source,omitempty"`
}

// Do 一键兑换主流币交易（POST /api/v5/trade/easy-convert）。
func (s *EasyConvertService) Do(ctx context.Context) ([]EasyConvertAck, error) {
	if len(s.fromCcy) == 0 || s.toCcy == "" {
		return nil, errEasyConvertMissingRequired
	}
	if len(s.fromCcy) > 5 {
		return nil, errEasyConvertTooManyFromCcy
	}
	for _, ccy := range s.fromCcy {
		if ccy == "" {
			return nil, errEasyConvertEmptyFromCcy
		}
		if ccy == s.toCcy {
			return nil, errEasyConvertSameCurrency
		}
	}

	req := easyConvertRequest{
		FromCcy: s.fromCcy,
		ToCcy:   s.toCcy,
		Source:  s.source,
	}

	var data []EasyConvertAck
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/trade/easy-convert", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyEasyConvertResponse
	}
	return data, nil
}
