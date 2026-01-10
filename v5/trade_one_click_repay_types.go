package okx

// OneClickRepayDebtData 表示一键还债负债币种信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type OneClickRepayDebtData struct {
	DebtAmt string `json:"debtAmt"`
	DebtCcy string `json:"debtCcy"`
}

// OneClickRepayRepayData 表示一键还债偿还币种信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type OneClickRepayRepayData struct {
	RepayAmt string `json:"repayAmt"`
	RepayCcy string `json:"repayCcy"`
}

// OneClickRepayCurrencyList 表示一键还债币种列表返回项（跨币种保证金/组合保证金）。
type OneClickRepayCurrencyList struct {
	DebtData  []OneClickRepayDebtData  `json:"debtData"`
	DebtType  string                   `json:"debtType"`
	RepayData []OneClickRepayRepayData `json:"repayData"`
}

// OneClickRepayAck 表示一键还债交易返回项（跨币种保证金/组合保证金）。
// 数值字段按 OKX 返回保持为 string（无损）。
type OneClickRepayAck struct {
	DebtCcy  string `json:"debtCcy"`
	RepayCcy string `json:"repayCcy"`
	Status   string `json:"status"`

	FillDebtSz  string `json:"fillDebtSz"`
	FillRepaySz string `json:"fillRepaySz"`

	FillFromSz string `json:"fillFromSz"`
	FillToSz   string `json:"fillToSz"`

	UTime int64 `json:"uTime,string"`
}

// OneClickRepayHistory 表示一键还债历史记录返回项（跨币种保证金/组合保证金）。
// 数值字段按 OKX 返回保持为 string（无损）。
type OneClickRepayHistory struct {
	DebtCcy     string `json:"debtCcy"`
	FillDebtSz  string `json:"fillDebtSz"`
	FillRepaySz string `json:"fillRepaySz"`
	RepayCcy    string `json:"repayCcy"`
	Status      string `json:"status"`
	UTime       int64  `json:"uTime,string"`
}
