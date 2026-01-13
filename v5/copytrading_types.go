package okx

// CopyTradingConfig 表示跟单/带单相关的账户配置信息。
// 数值字段保持为 string（无损）。
type CopyTradingConfig struct {
	UniqueCode string                    `json:"uniqueCode"`
	NickName   string                    `json:"nickName"`
	PortLink   string                    `json:"portLink"`
	Details    []CopyTradingConfigDetail `json:"details"`
}

// CopyTradingConfigDetail 表示账户配置详情（按 instType/roleType 等维度）。
type CopyTradingConfigDetail struct {
	CopyTraderNum      string `json:"copyTraderNum"`
	InstType           string `json:"instType"`
	MaxCopyTraderNum   string `json:"maxCopyTraderNum"`
	ProfitSharingRatio string `json:"profitSharingRatio"`
	RoleType           string `json:"roleType"`
}

// CopyTradingCopySettings 表示针对某交易员的跟单设置。
// 数值字段保持为 string（无损）。
type CopyTradingCopySettings struct {
	Ccy             string                        `json:"ccy"`
	CopyAmt         string                        `json:"copyAmt"`
	CopyInstIdType  string                        `json:"copyInstIdType"`
	CopyMgnMode     string                        `json:"copyMgnMode"`
	CopyMode        string                        `json:"copyMode"`
	CopyRatio       string                        `json:"copyRatio"`
	CopyState       string                        `json:"copyState"`
	CopyTotalAmt    string                        `json:"copyTotalAmt"`
	InstIds         []CopyTradingCopySettingsInst `json:"instIds"`
	SlRatio         string                        `json:"slRatio"`
	SlTotalAmt      string                        `json:"slTotalAmt"`
	SubPosCloseType string                        `json:"subPosCloseType"`
	TpRatio         string                        `json:"tpRatio"`
	Tag             string                        `json:"tag"`
}

// CopyTradingCopySettingsInst 表示跟单合约列表条目。
type CopyTradingCopySettingsInst struct {
	Enabled string `json:"enabled"`
	InstId  string `json:"instId"`
}

// CopyTradingCurrentLeadTrader 表示当前跟随的交易员信息。
// 数值字段保持为 string（无损）。
type CopyTradingCurrentLeadTrader struct {
	BeginCopyTime      UnixMilli `json:"beginCopyTime"`
	Ccy                string    `json:"ccy"`
	CopyTotalAmt       string    `json:"copyTotalAmt"`
	CopyTotalPnl       string    `json:"copyTotalPnl"`
	LeadMode           string    `json:"leadMode"`
	Margin             string    `json:"margin"`
	NickName           string    `json:"nickName"`
	PortLink           string    `json:"portLink"`
	ProfitSharingRatio string    `json:"profitSharingRatio"`
	TodayPnl           string    `json:"todayPnl"`
	UniqueCode         string    `json:"uniqueCode"`
	Upl                string    `json:"upl"`
}

// CopyTradingSubPosition 表示带单/跟单仓位信息（当前/历史/公开查询均复用）。
// 数值字段保持为 string（无损）。
type CopyTradingSubPosition struct {
	AlgoId           string    `json:"algoId"`
	AvailSubPos      string    `json:"availSubPos"`
	Ccy              string    `json:"ccy"`
	CloseAvgPx       string    `json:"closeAvgPx"`
	CloseSubPos      string    `json:"closeSubPos"`
	CloseTime        UnixMilli `json:"closeTime"`
	InstId           string    `json:"instId"`
	InstType         string    `json:"instType"`
	Lever            string    `json:"lever"`
	Margin           string    `json:"margin"`
	MarkPx           string    `json:"markPx"`
	MgnMode          string    `json:"mgnMode"`
	OpenAvgPx        string    `json:"openAvgPx"`
	OpenOrdId        string    `json:"openOrdId"`
	OpenTime         UnixMilli `json:"openTime"`
	Pnl              string    `json:"pnl"`
	PnlRatio         string    `json:"pnlRatio"`
	PosSide          string    `json:"posSide"`
	ProfitSharingAmt string    `json:"profitSharingAmt"`
	SlOrdPx          string    `json:"slOrdPx"`
	SlTriggerPx      string    `json:"slTriggerPx"`
	SubPos           string    `json:"subPos"`
	SubPosId         string    `json:"subPosId"`
	TpOrdPx          string    `json:"tpOrdPx"`
	TpTriggerPx      string    `json:"tpTriggerPx"`
	Type             string    `json:"type"`
	UniqueCode       string    `json:"uniqueCode"`
	Upl              string    `json:"upl"`
	UplRatio         string    `json:"uplRatio"`
}

// CopyTradingInstrument 表示带单产品设置。
type CopyTradingInstrument struct {
	Enabled bool   `json:"enabled"`
	InstId  string `json:"instId"`
}

// CopyTradingProfitSharingDetail 表示交易员历史分润明细。
// 数值字段保持为 string（无损）。
type CopyTradingProfitSharingDetail struct {
	Ccy              string    `json:"ccy"`
	InstType         string    `json:"instType"`
	NickName         string    `json:"nickName"`
	PortLink         string    `json:"portLink"`
	ProfitSharingAmt string    `json:"profitSharingAmt"`
	ProfitSharingId  string    `json:"profitSharingId"`
	TS               UnixMilli `json:"ts"`
}

// CopyTradingTotalProfitSharing 表示交易员历史分润汇总。
// 数值字段保持为 string（无损）。
type CopyTradingTotalProfitSharing struct {
	Ccy                   string `json:"ccy"`
	InstType              string `json:"instType"`
	TotalProfitSharingAmt string `json:"totalProfitSharingAmt"`
}

// CopyTradingUnrealizedProfitSharingDetail 表示交易员待分润明细。
// 数值字段保持为 string（无损）。
type CopyTradingUnrealizedProfitSharingDetail struct {
	Ccy                        string    `json:"ccy"`
	InstType                   string    `json:"instType"`
	NickName                   string    `json:"nickName"`
	PortLink                   string    `json:"portLink"`
	TS                         UnixMilli `json:"ts"`
	UnrealizedProfitSharingAmt string    `json:"unrealizedProfitSharingAmt"`
}

// CopyTradingTotalUnrealizedProfitSharing 表示交易员待分润汇总。
// 数值字段保持为 string（无损）。
type CopyTradingTotalUnrealizedProfitSharing struct {
	ProfitSharingTs                 UnixMilli `json:"profitSharingTs"`
	TotalUnrealizedProfitSharingAmt string    `json:"totalUnrealizedProfitSharingAmt"`
}

// CopyTradingResult 表示 result=true/false 的通用返回。
type CopyTradingResult struct {
	Result bool `json:"result"`
}

// CopyTradingSubPositionAck 表示带单/跟单仓位操作返回项（如设置止盈止损/平仓）。
type CopyTradingSubPositionAck struct {
	SubPosId string `json:"subPosId"`
	Tag      string `json:"tag"`
}

// CopyTradingPublicConfig 表示跟单设置的参数配置信息（公共）。
// 数值字段保持为 string（无损）。
type CopyTradingPublicConfig struct {
	MaxCopyAmt      string `json:"maxCopyAmt"`
	MinCopyAmt      string `json:"minCopyAmt"`
	MaxCopyTotalAmt string `json:"maxCopyTotalAmt"`
	MinCopyRatio    string `json:"minCopyRatio"`
	MaxCopyRatio    string `json:"maxCopyRatio"`
	MaxTpRatio      string `json:"maxTpRatio"`
	MaxSlRatio      string `json:"maxSlRatio"`
}

// CopyTradingPublicLeadTraders 表示交易员排名（公共）。
type CopyTradingPublicLeadTraders struct {
	DataVer   string                      `json:"dataVer"`
	TotalPage string                      `json:"totalPage"`
	Ranks     []CopyTradingLeadTraderRank `json:"ranks"`
}

// CopyTradingLeadTraderRank 表示交易员排名项。
// 数值字段保持为 string（无损）。
type CopyTradingLeadTraderRank struct {
	Aum              string                `json:"aum"`
	AccCopyTraderNum string                `json:"accCopyTraderNum"`
	Ccy              string                `json:"ccy"`
	CopyState        string                `json:"copyState"`
	CopyTraderNum    string                `json:"copyTraderNum"`
	LeadDays         string                `json:"leadDays"`
	MaxCopyTraderNum string                `json:"maxCopyTraderNum"`
	NickName         string                `json:"nickName"`
	Pnl              string                `json:"pnl"`
	PnlRatio         string                `json:"pnlRatio"`
	PnlRatios        []CopyTradingPnlRatio `json:"pnlRatios"`
	PortLink         string                `json:"portLink"`
	TraderInsts      []string              `json:"traderInsts"`
	UniqueCode       string                `json:"uniqueCode"`
	WinRatio         string                `json:"winRatio"`
}

// CopyTradingPnlRatio 表示某日收益率数据点。
type CopyTradingPnlRatio struct {
	BeginTs  UnixMilli `json:"beginTs"`
	PnlRatio string    `json:"pnlRatio"`
}

// CopyTradingPnl 表示收益表现数据项（周/日）。
// 数值字段保持为 string（无损）。
type CopyTradingPnl struct {
	BeginTs  UnixMilli `json:"beginTs"`
	Pnl      string    `json:"pnl"`
	PnlRatio string    `json:"pnlRatio"`
}

// CopyTradingPublicStats 表示交易员带单情况（公共）。
// 数值字段保持为 string（无损）。
type CopyTradingPublicStats struct {
	AvgSubPosNotional string `json:"avgSubPosNotional"`
	Ccy               string `json:"ccy"`
	CurCopyTraderPnl  string `json:"curCopyTraderPnl"`
	InvestAmt         string `json:"investAmt"`
	LossDays          string `json:"lossDays"`
	ProfitDays        string `json:"profitDays"`
	WinRatio          string `json:"winRatio"`
}

// CopyTradingPreferenceCurrency 表示交易员币种偏好（公共）。
// 数值字段保持为 string（无损）。
type CopyTradingPreferenceCurrency struct {
	Ccy   string `json:"ccy"`
	Ratio string `json:"ratio"`
}

// CopyTradingPublicCopyTraders 表示交易员跟单人信息（公共）。
// 数值字段保持为 string（无损）。
type CopyTradingPublicCopyTraders struct {
	Ccy                   string                  `json:"ccy"`
	CopyTotalPnl          string                  `json:"copyTotalPnl"`
	CopyTraderNumChg      string                  `json:"copyTraderNumChg"`
	CopyTraderNumChgRatio string                  `json:"copyTraderNumChgRatio"`
	CopyTraders           []CopyTradingCopyTrader `json:"copyTraders"`
}

// CopyTradingCopyTrader 表示跟单员信息。
// 数值字段保持为 string（无损）。
type CopyTradingCopyTrader struct {
	BeginCopyTime UnixMilli `json:"beginCopyTime"`
	NickName      string    `json:"nickName"`
	Pnl           string    `json:"pnl"`
	PortLink      string    `json:"portLink"`
}
