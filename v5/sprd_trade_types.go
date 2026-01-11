package okx

// SprdTrade 表示价差交易成交明细（REST/WS 共用的精简结构）。
//
// 说明：价格/数量等字段保持为 string（无损）；时间戳字段解析为 int64（json string）。
type SprdTrade struct {
	SprdId string `json:"sprdId"`

	TradeId string `json:"tradeId"`
	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	Tag     string `json:"tag"`

	FillPx string `json:"fillPx"`
	FillSz string `json:"fillSz"`

	State    string `json:"state"`
	Side     string `json:"side"`
	ExecType string `json:"execType"`

	TS int64 `json:"ts,string"`

	Legs []SprdTradeLeg `json:"legs"`

	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type SprdTradeLeg struct {
	InstId string `json:"instId"`

	Px     string `json:"px"`
	Sz     string `json:"sz"`
	SzCont string `json:"szCont"`
	Side   string `json:"side"`

	FillPnl string `json:"fillPnl"`
	Fee     string `json:"fee"`
	FeeCcy  string `json:"feeCcy"`

	TradeId string `json:"tradeId"`
}

// 兼容：历史命名（WS 推送成交）。
type WSSprdTrade = SprdTrade
type WSSprdTradeLeg = SprdTradeLeg
