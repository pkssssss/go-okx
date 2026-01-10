package okx

// OneClickRepayCurrencyListV2Item 表示一键还债币种列表（新）。
type OneClickRepayCurrencyListV2Item struct {
	DebtData  []OneClickRepayDebtData  `json:"debtData"`
	RepayData []OneClickRepayRepayData `json:"repayData"`

	DebtCcy      string   `json:"debtCcy"`
	RepayCcyList []string `json:"repayCcyList"`
	TS           int64    `json:"ts,string"`
}

// OneClickRepayV2Ack 表示一键还债交易（新）返回项。
type OneClickRepayV2Ack struct {
	DebtCcy      string   `json:"debtCcy"`
	RepayCcyList []string `json:"repayCcyList"`
	TS           int64    `json:"ts,string"`
}

// OneClickRepayHistoryV2OrderInfo 表示一键还债历史相关订单信息。
type OneClickRepayHistoryV2OrderInfo struct {
	OrdId   string `json:"ordId"`
	InstId  string `json:"instId"`
	OrdType string `json:"ordType"`
	Side    string `json:"side"`
	Px      string `json:"px"`
	Sz      string `json:"sz"`

	FillPx string `json:"fillPx"`
	FillSz string `json:"fillSz"`

	State string `json:"state"`
	CTime int64  `json:"cTime,string"`
}

// OneClickRepayHistoryV2Item 表示一键还债历史记录（新）。
type OneClickRepayHistoryV2Item struct {
	DebtCcy      string                            `json:"debtCcy"`
	FillDebtSz   string                            `json:"fillDebtSz"`
	OrdIdInfo    []OneClickRepayHistoryV2OrderInfo `json:"ordIdInfo"`
	RepayCcyList []string                          `json:"repayCcyList"`
	Status       string                            `json:"status"`
	TS           int64                             `json:"ts,string"`
}
