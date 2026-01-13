package okx

// ===== Flexible Loan（活期借币）=====

// FinanceFlexibleLoanBorrowCurrency 表示可借币种项。
type FinanceFlexibleLoanBorrowCurrency struct {
	BorrowCcy string `json:"borrowCcy"`
}

// FinanceFlexibleLoanCollateralAsset 表示可抵押资产明细。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanCollateralAsset struct {
	Amt         string `json:"amt"`
	Ccy         string `json:"ccy"`
	NotionalUsd string `json:"notionalUsd"`
}

// FinanceFlexibleLoanCollateralAssets 表示可抵押资产返回项。
type FinanceFlexibleLoanCollateralAssets struct {
	Assets []FinanceFlexibleLoanCollateralAsset `json:"assets"`
}

// FinanceFlexibleLoanSupCollateral 表示最大可借请求中的补充抵押资产信息。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanSupCollateral struct {
	Ccy string `json:"ccy"`
	Amt string `json:"amt"`
}

// FinanceFlexibleLoanMaxLoan 表示最大可借返回项。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanMaxLoan struct {
	BorrowCcy       string `json:"borrowCcy"`
	MaxLoan         string `json:"maxLoan"`
	NotionalUsd     string `json:"notionalUsd"`
	RemainingQuota  string `json:"remainingQuota"`
	AvailLoan       string `json:"availLoan,omitempty"`
	LoanNotionalUsd string `json:"loanNotionalUsd,omitempty"`
}

// FinanceFlexibleLoanInfoCcyAmt 表示借贷/抵押币种与数量对。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanInfoCcyAmt struct {
	Ccy string `json:"ccy"`
	Amt string `json:"amt"`
}

// FinanceFlexibleLoanRiskWarningData 表示风险预警信息。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanRiskWarningData struct {
	InstId string `json:"instId"`
	LiqPx  string `json:"liqPx"`
}

// FinanceFlexibleLoanInfo 表示借贷信息。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanInfo struct {
	LoanNotionalUsd       string `json:"loanNotionalUsd"`
	CollateralNotionalUsd string `json:"collateralNotionalUsd"`

	LoanData       []FinanceFlexibleLoanInfoCcyAmt `json:"loanData"`
	CollateralData []FinanceFlexibleLoanInfoCcyAmt `json:"collateralData"`

	RiskWarningData FinanceFlexibleLoanRiskWarningData `json:"riskWarningData"`

	CurLTV        string `json:"curLTV"`
	MarginCallLTV string `json:"marginCallLTV"`
	LiqLTV        string `json:"liqLTV"`
}

// FinanceFlexibleLoanLoanHistory 表示借贷历史条目。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanLoanHistory struct {
	RefId string    `json:"refId"`
	Type  string    `json:"type"`
	Ccy   string    `json:"ccy"`
	Amt   string    `json:"amt"`
	TS    UnixMilli `json:"ts"`
}

// FinanceFlexibleLoanInterestAccrued 表示计息记录条目。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanInterestAccrued struct {
	RefId        string    `json:"refId"`
	Ccy          string    `json:"ccy"`
	Loan         string    `json:"loan"`
	Interest     string    `json:"interest"`
	InterestRate string    `json:"interestRate"`
	TS           UnixMilli `json:"ts"`
}

// FinanceFlexibleLoanMaxCollateralRedeemAmount 表示最大可赎回抵押物数量。
// 数值字段保持为 string（无损）。
type FinanceFlexibleLoanMaxCollateralRedeemAmount struct {
	Ccy          string `json:"ccy"`
	MaxRedeemAmt string `json:"maxRedeemAmt"`
}

// ===== Savings（活期简单赚币）=====

// FinanceSavingsBalance 表示活期简单赚币余额条目。
// 数值字段保持为 string（无损）。
type FinanceSavingsBalance struct {
	Ccy        string `json:"ccy"`
	Amt        string `json:"amt"`
	Earnings   string `json:"earnings"`
	Rate       string `json:"rate"`
	LoanAmt    string `json:"loanAmt"`
	PendingAmt string `json:"pendingAmt"`
	RedemptAmt string `json:"redemptAmt"`
}

// FinanceSavingsPurchaseRedemptAck 表示活期简单赚币申购/赎回返回项。
// 数值字段保持为 string（无损）。
type FinanceSavingsPurchaseRedemptAck struct {
	Ccy  string `json:"ccy"`
	Amt  string `json:"amt"`
	Side string `json:"side"`
	Rate string `json:"rate"`
}

// FinanceSavingsSetLendingRateAck 表示设置活期简单赚币借贷利率返回项。
// 数值字段保持为 string（无损）。
type FinanceSavingsSetLendingRateAck struct {
	Ccy  string `json:"ccy"`
	Rate string `json:"rate"`
}

// FinanceSavingsLendingHistory 表示活期简单赚币出借明细条目。
// 数值字段保持为 string（无损）。
type FinanceSavingsLendingHistory struct {
	Ccy      string    `json:"ccy"`
	Amt      string    `json:"amt"`
	Earnings string    `json:"earnings"`
	Rate     string    `json:"rate"`
	TS       UnixMilli `json:"ts"`
}

// FinanceSavingsLendingRateSummary 表示市场借贷信息（公共）。
// 数值字段保持为 string（无损）。
type FinanceSavingsLendingRateSummary struct {
	Ccy       string `json:"ccy"`
	AvgAmt    string `json:"avgAmt"`
	AvgAmtUsd string `json:"avgAmtUsd"`
	AvgRate   string `json:"avgRate"`
	PreRate   string `json:"preRate"`
	EstRate   string `json:"estRate"`
}

// FinanceSavingsLendingRateHistory 表示市场借贷历史（公共）。
// 数值字段保持为 string（无损）。
type FinanceSavingsLendingRateHistory struct {
	Ccy  string    `json:"ccy"`
	Amt  string    `json:"amt"`
	Rate string    `json:"rate"`
	TS   UnixMilli `json:"ts"`
}

// ===== Staking / DeFi（链上赚币 / ETH / SOL 质押）=====

// FinanceStakingDefiOfferInvestData 表示项目可投币种信息。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiOfferInvestData struct {
	Ccy    string `json:"ccy"`
	Bal    string `json:"bal"`
	MinAmt string `json:"minAmt"`
	MaxAmt string `json:"maxAmt"`
}

// FinanceStakingDefiInvestData 表示投资信息（下单/订单明细）。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiInvestData struct {
	Ccy string `json:"ccy"`
	Amt string `json:"amt"`
}

// FinanceStakingDefiEarningData 表示收益信息。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiEarningData struct {
	Ccy              string `json:"ccy"`
	EarningType      string `json:"earningType"`
	Earnings         string `json:"earnings,omitempty"`
	RealizedEarnings string `json:"realizedEarnings,omitempty"`
}

// FinanceStakingDefiFastRedemptionData 表示快速赎回信息。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiFastRedemptionData struct {
	Ccy          string `json:"ccy"`
	RedeemingAmt string `json:"redeemingAmt"`
}

// FinanceStakingDefiOffer 表示链上赚币项目。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiOffer struct {
	ProductId    string `json:"productId"`
	Protocol     string `json:"protocol"`
	ProtocolType string `json:"protocolType"`
	Term         string `json:"term"`
	Apy          string `json:"apy"`
	EarlyRedeem  bool   `json:"earlyRedeem"`
	State        string `json:"state"`

	InvestData  []FinanceStakingDefiOfferInvestData `json:"investData"`
	EarningData []FinanceStakingDefiEarningData     `json:"earningData"`

	RedeemPeriod             []string `json:"redeemPeriod"`
	FastRedemptionDailyLimit string   `json:"fastRedemptionDailyLimit"`
}

// FinanceStakingDefiOrderAck 表示链上赚币下单/赎回/撤销返回项。
type FinanceStakingDefiOrderAck struct {
	OrdId string `json:"ordId"`
	Tag   string `json:"tag"`
}

// FinanceStakingDefiOrder 表示链上赚币订单（活跃/历史）。
// 数值字段保持为 string/UnixMilli（无损）。
type FinanceStakingDefiOrder struct {
	OrdId        string `json:"ordId"`
	Ccy          string `json:"ccy"`
	ProductId    string `json:"productId"`
	State        string `json:"state"`
	Protocol     string `json:"protocol"`
	ProtocolType string `json:"protocolType"`
	Term         string `json:"term"`
	Apy          string `json:"apy"`

	InvestData         []FinanceStakingDefiInvestData         `json:"investData"`
	EarningData        []FinanceStakingDefiEarningData        `json:"earningData"`
	FastRedemptionData []FinanceStakingDefiFastRedemptionData `json:"fastRedemptionData"`

	PurchasedTime            UnixMilli `json:"purchasedTime"`
	RedeemedTime             UnixMilli `json:"redeemedTime,omitempty"`
	EstSettlementTime        UnixMilli `json:"estSettlementTime"`
	CancelRedemptionDeadline UnixMilli `json:"cancelRedemptionDeadline"`

	Tag string `json:"tag"`
}

// FinanceStakingDefiPurchaseRedeemHistory 表示申购/赎回记录（ETH/SOL）。
// 数值字段保持为 string/UnixMilli（无损）。
type FinanceStakingDefiPurchaseRedeemHistory struct {
	Type             string    `json:"type"`
	Amt              string    `json:"amt"`
	RedeemingAmt     string    `json:"redeemingAmt"`
	Status           string    `json:"status"`
	OrdId            string    `json:"ordId"`
	RequestTime      UnixMilli `json:"requestTime"`
	CompletedTime    UnixMilli `json:"completedTime"`
	EstCompletedTime UnixMilli `json:"estCompletedTime"`
}

// FinanceStakingDefiAPYHistory 表示历史收益率（公共）。
// 数值字段保持为 string/UnixMilli（无损）。
type FinanceStakingDefiAPYHistory struct {
	Rate string    `json:"rate"`
	TS   UnixMilli `json:"ts"`
}

// FinanceStakingDefiETHProductInfo 表示 ETH 质押产品信息。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiETHProductInfo struct {
	FastRedemptionDailyLimit string `json:"fastRedemptionDailyLimit"`
}

// FinanceStakingDefiETHBalance 表示 ETH 质押余额快照。
// 数值字段保持为 string/UnixMilli（无损）。
type FinanceStakingDefiETHBalance struct {
	Amt                   string    `json:"amt"`
	Ccy                   string    `json:"ccy"`
	LatestInterestAccrual string    `json:"latestInterestAccrual"`
	TotalInterestAccrual  string    `json:"totalInterestAccrual"`
	TS                    UnixMilli `json:"ts"`
}

// FinanceStakingDefiSOLProductInfo 表示 SOL 质押产品信息。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiSOLProductInfo struct {
	FastRedemptionDailyLimit string `json:"fastRedemptionDailyLimit"`
	FastRedemptionAvail      string `json:"fastRedemptionAvail"`
}

// FinanceStakingDefiSOLBalance 表示 SOL 质押余额（OKSOL）。
// 数值字段保持为 string（无损）。
type FinanceStakingDefiSOLBalance struct {
	Amt                   string `json:"amt"`
	Ccy                   string `json:"ccy"`
	LatestInterestAccrual string `json:"latestInterestAccrual"`
	TotalInterestAccrual  string `json:"totalInterestAccrual"`
}
