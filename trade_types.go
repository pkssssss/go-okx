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
	PosSide string `json:"posSide"`
	TdMode  string `json:"tdMode"`
	OrdType string `json:"ordType"`
	State   string `json:"state"`

	Ccy           string `json:"ccy"`
	TgtCcy        string `json:"tgtCcy"`
	TradeQuoteCcy string `json:"tradeQuoteCcy"`
	ReduceOnly    string `json:"reduceOnly"`

	Px     string `json:"px"`
	PxUsd  string `json:"pxUsd"`
	PxVol  string `json:"pxVol"`
	PxType string `json:"pxType"`
	Sz     string `json:"sz"`

	AvgPx     string `json:"avgPx"`
	FillPx    string `json:"fillPx"`
	FillSz    string `json:"fillSz"`
	AccFillSz string `json:"accFillSz"`
	TradeId   string `json:"tradeId"`
	FillTime  string `json:"fillTime"`

	Pnl       string `json:"pnl"`
	Fee       string `json:"fee"`
	FeeCcy    string `json:"feeCcy"`
	Rebate    string `json:"rebate"`
	RebateCcy string `json:"rebateCcy"`

	StpMode            string `json:"stpMode"`
	CancelSource       string `json:"cancelSource"`
	CancelSourceReason string `json:"cancelSourceReason"`

	UTime int64 `json:"uTime,string"`
	CTime int64 `json:"cTime,string"`
}

// TradeFill 表示成交明细（精简版）。
// 价格/数量字段保持为 string（无损），时间戳字段解析为 int64。
type TradeFill struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	TradeId string `json:"tradeId"`
	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`

	BillId  string `json:"billId"`
	SubType string `json:"subType"`
	Tag     string `json:"tag"`

	FillPx    string `json:"fillPx"`
	FillSz    string `json:"fillSz"`
	FillIdxPx string `json:"fillIdxPx"`
	FillPnl   string `json:"fillPnl"`

	FillPxVol   string `json:"fillPxVol"`
	FillPxUsd   string `json:"fillPxUsd"`
	FillMarkVol string `json:"fillMarkVol"`
	FillFwdPx   string `json:"fillFwdPx"`
	FillMarkPx  string `json:"fillMarkPx"`

	Side     string `json:"side"`
	PosSide  string `json:"posSide"`
	ExecType string `json:"execType"`

	FeeCcy string `json:"feeCcy"`
	Fee    string `json:"fee"`

	TS       int64 `json:"ts,string"`
	FillTime int64 `json:"fillTime,string"`
}
