package okx

// TradingBotOrderAck 表示 TradingBot 下单/改单/撤单等通用返回项。
// 数值字段保持为 string（无损）。
type TradingBotOrderAck struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`

	SCode string `json:"sCode"`
	SMsg  string `json:"sMsg"`

	Tag string `json:"tag"`
}

// TradingBotGridAIParam 表示网格策略智能回测结果。
// 数值字段保持为 string（无损）。
type TradingBotGridAIParam struct {
	InstId      string `json:"instId"`
	AlgoOrdType string `json:"algoOrdType"`
	Duration    string `json:"duration"`
	GridNum     string `json:"gridNum"`
	MaxPx       string `json:"maxPx"`
	MinPx       string `json:"minPx"`

	PerMaxProfitRate   string `json:"perMaxProfitRate"`
	PerMinProfitRate   string `json:"perMinProfitRate"`
	PerGridProfitRatio string `json:"perGridProfitRatio"`
	AnnualizedRate     string `json:"annualizedRate"`

	MinInvestment string `json:"minInvestment"`
	Ccy           string `json:"ccy"`
	RunType       string `json:"runType"`

	Direction string `json:"direction"`
	Lever     string `json:"lever"`
	SourceCcy string `json:"sourceCcy"`
}

// TradingBotGridMaxGridQty 表示最大网格数量返回项。
type TradingBotGridMaxGridQty struct {
	MaxGridQty string `json:"maxGridQty"`
}

// TradingBotPublicRSIBackTestingResult 表示 RSI 回测结果。
type TradingBotPublicRSIBackTestingResult struct {
	TriggerNum string `json:"triggerNum"`
}

// TradingBotGridTriggerParam 表示网格信号触发参数（下单/查询复用）。
type TradingBotGridTriggerParam struct {
	TriggerAction   string `json:"triggerAction"`
	TriggerStrategy string `json:"triggerStrategy"`

	DelaySeconds string `json:"delaySeconds,omitempty"`
	Timeframe    string `json:"timeframe,omitempty"`
	Thold        string `json:"thold,omitempty"`
	TriggerCond  string `json:"triggerCond,omitempty"`
	TimePeriod   string `json:"timePeriod,omitempty"`

	TriggerPx string `json:"triggerPx,omitempty"`
	StopType  string `json:"stopType,omitempty"`

	TriggerTime string `json:"triggerTime,omitempty"`
	TriggerType string `json:"triggerType,omitempty"`
}

// TradingBotGridRebateTrans 表示返佣划转信息。
type TradingBotGridRebateTrans struct {
	Rebate    string `json:"rebate"`
	RebateCcy string `json:"rebateCcy"`
}

// TradingBotGridOrder 表示网格策略委托单信息（详情/未完成/历史复用）。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotGridOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	Uly      string `json:"uly"`

	AlgoOrdType string `json:"algoOrdType"`
	State       string `json:"state"`
	Tag         string `json:"tag"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	RunType   string `json:"runType"`
	Direction string `json:"direction"`

	Lever       string `json:"lever"`
	ActualLever string `json:"actualLever"`

	MaxPx   string `json:"maxPx"`
	MinPx   string `json:"minPx"`
	GridNum string `json:"gridNum"`

	QuoteSz string `json:"quoteSz"`
	BaseSz  string `json:"baseSz"`
	Sz      string `json:"sz"`

	TradeQuoteCcy string `json:"tradeQuoteCcy"`

	BasePos bool `json:"basePos"`

	TpTriggerPx string `json:"tpTriggerPx"`
	SlTriggerPx string `json:"slTriggerPx"`
	TpRatio     string `json:"tpRatio"`
	SlRatio     string `json:"slRatio"`

	TotalPnl       string `json:"totalPnl"`
	PnlRatio       string `json:"pnlRatio"`
	FloatPnl       string `json:"floatPnl"`
	AnnualizedRate string `json:"annualizedRate"`

	ActiveOrdNum string `json:"activeOrdNum"`
	ArbitrageNum string `json:"arbitrageNum"`

	Eq       string `json:"eq"`
	AvailEq  string `json:"availEq"`
	FrozenEq string `json:"frozenEq"`

	CurQuoteSz string `json:"curQuoteSz"`
	CurBaseSz  string `json:"curBaseSz"`

	CancelType string `json:"cancelType"`
	StopType   string `json:"stopType"`
	StopResult string `json:"stopResult"`

	ProfitSharingRatio string `json:"profitSharingRatio"`
	CopyType           string `json:"copyType"`
	Fee                string `json:"fee"`
	FundingFee         string `json:"fundingFee"`

	RebateTrans   []TradingBotGridRebateTrans  `json:"rebateTrans"`
	TriggerParams []TradingBotGridTriggerParam `json:"triggerParams"`
}

// TradingBotGridPosition 表示网格/信号策略持仓信息（字段结构相同）。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotGridPosition struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	AvgPx    string `json:"avgPx"`
	Ccy      string `json:"ccy"`
	Lever    string `json:"lever"`
	LiqPx    string `json:"liqPx"`
	PosSide  string `json:"posSide"`
	Pos      string `json:"pos"`
	MgnMode  string `json:"mgnMode"`
	MgnRatio string `json:"mgnRatio"`
	Imr      string `json:"imr"`
	Mmr      string `json:"mmr"`
	Upl      string `json:"upl"`
	UplRatio string `json:"uplRatio"`
	Last     string `json:"last"`

	NotionalUsd string `json:"notionalUsd"`
	Adl         string `json:"adl"`
	MarkPx      string `json:"markPx"`
}

// TradingBotSignalPosition 表示信号策略持仓（与网格持仓字段一致）。
type TradingBotSignalPosition = TradingBotGridPosition

// TradingBotGridSubOrder 表示网格策略子订单信息。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotGridSubOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	GroupId string `json:"groupId"`
	OrdId   string `json:"ordId"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	TdMode string `json:"tdMode"`
	Ccy    string `json:"ccy"`

	OrdType string `json:"ordType"`
	State   string `json:"state"`
	Side    string `json:"side"`
	PosSide string `json:"posSide"`

	Px        string `json:"px"`
	Sz        string `json:"sz"`
	AccFillSz string `json:"accFillSz"`
	AvgPx     string `json:"avgPx"`

	Fee       string `json:"fee"`
	FeeCcy    string `json:"feeCcy"`
	Rebate    string `json:"rebate"`
	RebateCcy string `json:"rebateCcy"`

	Pnl   string `json:"pnl"`
	CtVal string `json:"ctVal"`
	Lever string `json:"lever"`

	Tag string `json:"tag"`
}

// TradingBotGridInvestmentData 表示网格投资信息（min-investment 入参/出参复用）。
type TradingBotGridInvestmentData struct {
	Amt string `json:"amt"`
	Ccy string `json:"ccy"`
}

// TradingBotGridMinInvestmentResult 表示计算最小投资数量返回项。
type TradingBotGridMinInvestmentResult struct {
	MinInvestmentData []TradingBotGridInvestmentData `json:"minInvestmentData"`
	SingleAmt         string                         `json:"singleAmt"`
}

// TradingBotGridComputeMarginBalanceResult 表示调整保证金计算返回项。
type TradingBotGridComputeMarginBalanceResult struct {
	MaxAmt string `json:"maxAmt"`
	Lever  string `json:"lever"`
}

// TradingBotGridAmendAlgoBasicParamResult 表示修改网格策略基本参数返回项。
type TradingBotGridAmendAlgoBasicParamResult struct {
	AlgoId              string `json:"algoId"`
	RequiredTopupAmount string `json:"requiredTopupAmount"`
}

// TradingBotGridWithdrawIncomeAck 表示现货网格提取利润返回项。
type TradingBotGridWithdrawIncomeAck struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	Profit      string `json:"profit"`
}

// TradingBotAlgoIdAck 表示仅返回 algoId 的通用返回项。
type TradingBotAlgoIdAck struct {
	AlgoId string `json:"algoId"`
}

// TradingBotGridCloseOrderAck 表示合约网格平仓/撤销平仓单返回项。
type TradingBotGridCloseOrderAck struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	OrdId       string `json:"ordId"`
	Tag         string `json:"tag"`
}

// TradingBotRecurringListItem 表示定投币种配置/统计信息。
// 数值字段保持为 string（无损）。
type TradingBotRecurringListItem struct {
	Ccy   string `json:"ccy"`
	Ratio string `json:"ratio"`

	TotalAmt string `json:"totalAmt"`
	Profit   string `json:"profit"`
	AvgPx    string `json:"avgPx"`
	Px       string `json:"px"`
}

// TradingBotRecurringOrder 表示定投策略委托单信息（详情/未完成/历史复用）。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotRecurringOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`

	AlgoOrdType string `json:"algoOrdType"`
	InstType    string `json:"instType"`
	State       string `json:"state"`

	StgyName string `json:"stgyName"`

	Period        string `json:"period"`
	RecurringDay  string `json:"recurringDay"`
	RecurringHour string `json:"recurringHour"`
	RecurringTime string `json:"recurringTime"`
	TimeZone      string `json:"timeZone"`

	Amt           string `json:"amt"`
	InvestmentAmt string `json:"investmentAmt"`
	InvestmentCcy string `json:"investmentCcy"`

	NextInvestTime UnixMilli `json:"nextInvestTime"`

	TotalPnl     string `json:"totalPnl"`
	TotalAnnRate string `json:"totalAnnRate"`
	PnlRatio     string `json:"pnlRatio"`
	MktCap       string `json:"mktCap"`
	Cycles       string `json:"cycles"`

	Tag           string `json:"tag"`
	TradeQuoteCcy string `json:"tradeQuoteCcy"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	RecurringList []TradingBotRecurringListItem `json:"recurringList"`
}

// TradingBotRecurringSubOrder 表示定投策略子订单信息。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotRecurringSubOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OrdId   string `json:"ordId"`
	OrdType string `json:"ordType"`
	State   string `json:"state"`
	Side    string `json:"side"`

	Px        string `json:"px"`
	Sz        string `json:"sz"`
	AccFillSz string `json:"accFillSz"`
	AvgPx     string `json:"avgPx"`

	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	Tag    string `json:"tag"`
	TdMode string `json:"tdMode"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`
}

// TradingBotSignal 表示信号信息（signals 返回项）。
type TradingBotSignal struct {
	SignalChanId     string `json:"signalChanId"`
	SignalChanName   string `json:"signalChanName"`
	SignalChanDesc   string `json:"signalChanDesc"`
	SignalChanToken  string `json:"signalChanToken"`
	SignalSourceType string `json:"signalSourceType"`
}

// TradingBotSignalCreateAck 表示创建信号返回项。
type TradingBotSignalCreateAck struct {
	SignalChanId    string `json:"signalChanId"`
	SignalChanToken string `json:"signalChanToken"`
}

// TradingBotSignalEntrySettingParam 表示信号策略进场参数设定（下单/查询复用）。
type TradingBotSignalEntrySettingParam struct {
	AllowMultipleEntry *bool  `json:"allowMultipleEntry,omitempty"`
	EntryType          string `json:"entryType,omitempty"`
	Amt                string `json:"amt,omitempty"`
	Ratio              string `json:"ratio,omitempty"`
}

// TradingBotSignalExitSettingParam 表示信号策略离场参数设定（下单/查询复用）。
type TradingBotSignalExitSettingParam struct {
	TpSlType string `json:"tpSlType"`
	TpPct    string `json:"tpPct,omitempty"`
	SlPct    string `json:"slPct,omitempty"`
}

// TradingBotSignalOrder 表示信号策略信息（详情/活跃/历史复用）。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotSignalOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string   `json:"instType"`
	InstIds  []string `json:"instIds"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`

	State      string `json:"state"`
	CancelType string `json:"cancelType"`

	Lever      string `json:"lever"`
	InvestAmt  string `json:"investAmt"`
	SubOrdType string `json:"subOrdType"`
	Ratio      string `json:"ratio"`

	AvailBal  string `json:"availBal"`
	FrozenBal string `json:"frozenBal"`
	TotalEq   string `json:"totalEq"`

	TotalPnl      string `json:"totalPnl"`
	TotalPnlRatio string `json:"totalPnlRatio"`
	FloatPnl      string `json:"floatPnl"`
	RealizedPnl   string `json:"realizedPnl"`

	SignalChanId     string `json:"signalChanId"`
	SignalChanName   string `json:"signalChanName"`
	SignalSourceType string `json:"signalSourceType"`

	EntrySettingParam TradingBotSignalEntrySettingParam `json:"entrySettingParam"`
	ExitSettingParam  TradingBotSignalExitSettingParam  `json:"exitSettingParam"`
}

// TradingBotSignalPositionsHistory 表示信号策略历史仓位信息。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotSignalPositionsHistory struct {
	InstId     string    `json:"instId"`
	MgnMode    string    `json:"mgnMode"`
	CTime      UnixMilli `json:"cTime"`
	UTime      UnixMilli `json:"uTime"`
	OpenAvgPx  string    `json:"openAvgPx"`
	CloseAvgPx string    `json:"closeAvgPx"`
	Pnl        string    `json:"pnl"`
	PnlRatio   string    `json:"pnlRatio"`
	Lever      string    `json:"lever"`
	Direction  string    `json:"direction"`
	Uly        string    `json:"uly"`
}

// TradingBotSignalSubOrder 表示信号策略子订单信息。
// 数值字段保持为 string/UnixMilli（无损）。
type TradingBotSignalSubOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OrdId       string `json:"ordId"`
	SignalOrdId string `json:"signalOrdId"`
	ClOrdId     string `json:"clOrdId"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`
	PTime UnixMilli `json:"pTime"`

	TdMode string `json:"tdMode"`
	Ccy    string `json:"ccy"`

	OrdType string `json:"ordType"`
	State   string `json:"state"`
	Side    string `json:"side"`
	PosSide string `json:"posSide"`

	Px        string `json:"px"`
	Sz        string `json:"sz"`
	AccFillSz string `json:"accFillSz"`
	AvgPx     string `json:"avgPx"`

	Fee       string `json:"fee"`
	FeeCcy    string `json:"feeCcy"`
	Rebate    string `json:"rebate"`
	RebateCcy string `json:"rebateCcy"`

	Pnl   string `json:"pnl"`
	CtVal string `json:"ctVal"`
	Lever string `json:"lever"`

	Tag string `json:"tag"`
}

// TradingBotSignalCancelSubOrderAck 表示信号策略撤单返回项。
type TradingBotSignalCancelSubOrderAck struct {
	SignalOrdId string `json:"signalOrdId"`
	SCode       string `json:"sCode"`
	SMsg        string `json:"sMsg"`
}

// TradingBotSignalEventTriggeredOrdData 表示信号触发子订单信息。
type TradingBotSignalEventTriggeredOrdData struct {
	ClOrdId string `json:"clOrdId"`
}

// TradingBotSignalEventHistory 表示信号策略历史事件。
type TradingBotSignalEventHistory struct {
	AlertMsg         string                                  `json:"alertMsg"`
	AlgoId           string                                  `json:"algoId"`
	EventType        string                                  `json:"eventType"`
	EventCtime       UnixMilli                               `json:"eventCtime"`
	EventUtime       UnixMilli                               `json:"eventUtime"`
	EventProcessMsg  string                                  `json:"eventProcessMsg"`
	EventStatus      string                                  `json:"eventStatus"`
	TriggeredOrdData []TradingBotSignalEventTriggeredOrdData `json:"triggeredOrdData"`
}
