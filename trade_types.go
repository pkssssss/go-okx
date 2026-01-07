package okx

// TradeOrderAck 是 OKX 交易接口常见的返回项（下单/撤单/改单等）。
type TradeOrderAck struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	ReqId   string `json:"reqId"`
	Tag     string `json:"tag"`
	TS      int64  `json:"ts,string"`

	SCode string `json:"sCode"`
	SMsg  string `json:"sMsg"`
}

// TradeOrder 表示订单信息（精简版）。
type TradeOrder struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	Tag     string `json:"tag"`

	Side    string `json:"side"`
	OrdType string `json:"ordType"`
	State   string `json:"state"`

	Px string `json:"px"`
	Sz string `json:"sz"`

	AvgPx     string `json:"avgPx"`
	FillPx    string `json:"fillPx"`
	FillSz    string `json:"fillSz"`
	AccFillSz string `json:"accFillSz"`

	UTime int64 `json:"uTime,string"`
	CTime int64 `json:"cTime,string"`
}

// TradeFill 表示成交明细（精简版）。
type TradeFill struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	TradeId string `json:"tradeId"`

	Side    string `json:"side"`
	PosSide string `json:"posSide"`

	FillPx string `json:"fillPx"`
	FillSz string `json:"fillSz"`
	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	FillTime int64 `json:"fillTime,string"`
}
