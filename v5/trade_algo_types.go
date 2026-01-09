package okx

// TradeAlgoOrderAck 表示策略委托相关接口常见的返回项（下单/撤单/改单等）。
type TradeAlgoOrderAck struct {
	AlgoId      string `json:"algoId"`
	ClOrdId     string `json:"clOrdId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	ReqId       string `json:"reqId"`
	Tag         string `json:"tag"`

	SCode string `json:"sCode"`
	SMsg  string `json:"sMsg"`
}

// TradeAlgoOrderAttach 表示附带止盈止损信息（下单/查询返回）。
type TradeAlgoOrderAttach struct {
	AttachAlgoClOrdId string `json:"attachAlgoClOrdId,omitempty"`

	TpTriggerPx     string `json:"tpTriggerPx,omitempty"`
	TpTriggerPxType string `json:"tpTriggerPxType,omitempty"`
	TpOrdPx         string `json:"tpOrdPx,omitempty"`
	TpOrdKind       string `json:"tpOrdKind,omitempty"`

	SlTriggerPx     string `json:"slTriggerPx,omitempty"`
	SlTriggerPxType string `json:"slTriggerPxType,omitempty"`
	SlOrdPx         string `json:"slOrdPx,omitempty"`
	SlOrdKind       string `json:"slOrdKind,omitempty"`
}

// TradeAlgoOrderAttachAmend 表示修改附带止盈止损信息（改单请求）。
type TradeAlgoOrderAttachAmend struct {
	NewTpTriggerPx       string `json:"newTpTriggerPx,omitempty"`
	NewTpTriggerRatio    string `json:"newTpTriggerRatio,omitempty"`
	NewTpTriggerPxType   string `json:"newTpTriggerPxType,omitempty"`
	NewTpOrdPx           string `json:"newTpOrdPx,omitempty"`
	NewSlTriggerPx       string `json:"newSlTriggerPx,omitempty"`
	NewSlTriggerRatio    string `json:"newSlTriggerRatio,omitempty"`
	NewSlTriggerPxType   string `json:"newSlTriggerPxType,omitempty"`
	NewSlOrdPx           string `json:"newSlOrdPx,omitempty"`
	AttachAlgoClOrdId    string `json:"attachAlgoClOrdId,omitempty"`
	NewAttachAlgoClOrdId string `json:"newAttachAlgoClOrdId,omitempty"`
}

// TradeAlgoOrder 表示策略委托单信息（精简版；字段按 OKX 返回保持 string，时间戳字段解析为 int64）。
type TradeAlgoOrder struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	Ccy      string `json:"ccy"`

	OrdId     string   `json:"ordId"`
	OrdIdList []string `json:"ordIdList"`

	AlgoId      string `json:"algoId"`
	ClOrdId     string `json:"clOrdId"`
	AlgoClOrdId string `json:"algoClOrdId"`

	Sz            string `json:"sz"`
	CloseFraction string `json:"closeFraction"`

	OrdType string `json:"ordType"`
	Side    string `json:"side"`
	PosSide string `json:"posSide"`
	TdMode  string `json:"tdMode"`
	TgtCcy  string `json:"tgtCcy"`

	State string `json:"state"`

	Tag string `json:"tag"`

	TriggerPx     string `json:"triggerPx"`
	TriggerPxType string `json:"triggerPxType"`
	TriggerTime   string `json:"triggerTime"`
	OrderPx       string `json:"orderPx"`

	TpTriggerPx     string `json:"tpTriggerPx"`
	TpTriggerPxType string `json:"tpTriggerPxType"`
	TpOrdPx         string `json:"tpOrdPx"`
	SlTriggerPx     string `json:"slTriggerPx"`
	SlTriggerPxType string `json:"slTriggerPxType"`
	SlOrdPx         string `json:"slOrdPx"`

	CallbackRatio  string `json:"callbackRatio"`
	CallbackSpread string `json:"callbackSpread"`
	ActivePx       string `json:"activePx"`

	PxLimit      string `json:"pxLimit"`
	SzLimit      string `json:"szLimit"`
	TimeInterval string `json:"timeInterval"`
	PxSpread     string `json:"pxSpread"`

	TradeQuoteCcy string `json:"tradeQuoteCcy"`

	AttachAlgoOrds []TradeAlgoOrderAttach `json:"attachAlgoOrds"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`
}

// OrderPrecheckResult 表示订单预检查返回项（精简版；数值字段保持 string）。
type OrderPrecheckResult struct {
	AdjEq          string `json:"adjEq"`
	AdjEqChg       string `json:"adjEqChg"`
	Imr            string `json:"imr"`
	ImrChg         string `json:"imrChg"`
	Mmr            string `json:"mmr"`
	MmrChg         string `json:"mmrChg"`
	MgnRatio       string `json:"mgnRatio"`
	MgnRatioChg    string `json:"mgnRatioChg"`
	LiqPx          string `json:"liqPx"`
	LiqPxDiff      string `json:"liqPxDiff"`
	LiqPxDiffRatio string `json:"liqPxDiffRatio"`

	AvailBal    string `json:"availBal"`
	AvailBalChg string `json:"availBalChg"`
	Liab        string `json:"liab"`
	LiabChg     string `json:"liabChg"`
	LiabChgCcy  string `json:"liabChgCcy"`

	PosBal    string `json:"posBal"`
	PosBalChg string `json:"posBalChg"`

	Type string `json:"type"`
}
