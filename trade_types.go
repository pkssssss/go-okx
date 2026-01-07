package okx

// TradeOrderAck 是 OKX 交易接口常见的返回项（下单/撤单/改单等）。
type TradeOrderAck struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	Tag     string `json:"tag"`
	TS      string `json:"ts"`

	SCode string `json:"sCode"`
	SMsg  string `json:"sMsg"`
}
