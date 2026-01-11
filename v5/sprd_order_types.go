package okx

// SprdOrder 表示价差交易订单信息（REST/WS 共用的精简结构）。
//
// 说明：价格/数量等字段保持为 string（无损）；时间戳字段解析为 int64（json string）。
type SprdOrder struct {
	SprdId string `json:"sprdId"`

	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	Tag     string `json:"tag"`

	Px      string `json:"px"`
	Sz      string `json:"sz"`
	OrdType string `json:"ordType"`
	Side    string `json:"side"`

	FillSz  string `json:"fillSz"`
	FillPx  string `json:"fillPx"`
	TradeId string `json:"tradeId"`

	AccFillSz          string `json:"accFillSz"`
	PendingFillSz      string `json:"pendingFillSz"`
	PendingSettleSz    string `json:"pendingSettleSz"`
	CanceledSz         string `json:"canceledSz"`
	State              string `json:"state"`
	AvgPx              string `json:"avgPx"`
	CancelSource       string `json:"cancelSource"`
	CancelSourceReason string `json:"cancelSourceReason,omitempty"`

	Code        string `json:"code,omitempty"`
	Msg         string `json:"msg,omitempty"`
	ReqId       string `json:"reqId,omitempty"`
	AmendResult string `json:"amendResult,omitempty"`

	UTime int64 `json:"uTime,string"`
	CTime int64 `json:"cTime,string"`
}
