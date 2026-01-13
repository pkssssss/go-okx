package okx

import (
	"context"
	"net/http"
)

// FinanceStakingDefiSOLProductInfoService 获取 SOL 质押产品信息。
type FinanceStakingDefiSOLProductInfoService struct {
	c *Client
}

// NewFinanceStakingDefiSOLProductInfoService 创建 FinanceStakingDefiSOLProductInfoService。
func (c *Client) NewFinanceStakingDefiSOLProductInfoService() *FinanceStakingDefiSOLProductInfoService {
	return &FinanceStakingDefiSOLProductInfoService{c: c}
}

// Do 获取 SOL 质押产品信息（GET /api/v5/finance/staking-defi/sol/product-info）。
//
// 注意：OKX 文档示例中该接口的 data 为对象（非数组）。
func (s *FinanceStakingDefiSOLProductInfoService) Do(ctx context.Context) (*FinanceStakingDefiSOLProductInfo, error) {
	var data FinanceStakingDefiSOLProductInfo
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/finance/staking-defi/sol/product-info", nil, nil, true, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
