package okx

import (
	"encoding/json"
	"strings"
)

const (
	WSChannelOrders             = "orders"
	WSChannelFills              = "fills"
	WSChannelAccount            = "account"
	WSChannelPositions          = "positions"
	WSChannelBalanceAndPosition = "balance_and_position"

	WSChannelLiquidationWarning = "liquidation-warning"
	WSChannelAccountGreeks      = "account-greeks"

	WSChannelOrdersAlgo                  = "orders-algo"
	WSChannelAlgoAdvance                 = "algo-advance"
	WSChannelGridOrdersSpot              = "grid-orders-spot"
	WSChannelGridOrdersContract          = "grid-orders-contract"
	WSChannelGridPositions               = "grid-positions"
	WSChannelGridSubOrders               = "grid-sub-orders"
	WSChannelAlgoRecurringBuy            = "algo-recurring-buy"
	WSChannelCopytradingLeadNotification = "copytrading-lead-notification"

	WSChannelRFQs                   = "rfqs"
	WSChannelQuotes                 = "quotes"
	WSChannelStrucBlockTrades       = "struc-block-trades"
	WSChannelPublicStrucBlockTrades = "public-struc-block-trades"
	WSChannelPublicBlockTrades      = "public-block-trades"
	WSChannelBlockTickers           = "block-tickers"

	WSChannelDepositInfo    = "deposit-info"
	WSChannelWithdrawalInfo = "withdrawal-info"

	WSChannelEconomicCalendar = "economic-calendar"

	WSChannelInstruments = "instruments"
	WSChannelTickers     = "tickers"
	WSChannelTrades      = "trades"
	WSChannelTradesAll   = "trades-all"

	WSChannelStatus = "status"

	WSChannelOpenInterest = "open-interest"
	WSChannelFundingRate  = "funding-rate"
	WSChannelPriceLimit   = "price-limit"
	WSChannelMarkPrice    = "mark-price"
	WSChannelIndexTickers = "index-tickers"
	WSChannelOptSummary   = "opt-summary"

	WSChannelEstimatedPrice = "estimated-price"
	WSChannelADLWarning     = "adl-warning"

	WSChannelLiquidationOrders = "liquidation-orders"

	WSChannelBooks        = "books"
	WSChannelBooksELP     = "books-elp"
	WSChannelBooks5       = "books5"
	WSChannelBboTbt       = "bbo-tbt"
	WSChannelBooksL2Tbt   = "books-l2-tbt"
	WSChannelBooks50L2Tbt = "books50-l2-tbt"

	WSChannelSprdBboTbt     = "sprd-bbo-tbt"
	WSChannelSprdBooks5     = "sprd-books5"
	WSChannelSprdBooksL2Tbt = "sprd-books-l2-tbt"

	WSChannelOptionTrades       = "option-trades"
	WSChannelCallAuctionDetails = "call-auction-details"

	WSChannelSprdOrders = "sprd-orders"
	WSChannelSprdTrades = "sprd-trades"

	WSChannelSprdPublicTrades = "sprd-public-trades"
	WSChannelSprdTickers      = "sprd-tickers"
)

// WSEvent 表示 OKX WebSocket 的 event 消息（subscribe/login/error/notice 等）。
type WSEvent struct {
	ID    string `json:"id,omitempty"`
	Event string `json:"event"`
	Code  string `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`

	Arg    *WSArg `json:"arg,omitempty"`
	ConnID string `json:"connId,omitempty"`

	Channel   string `json:"channel,omitempty"`
	ConnCount string `json:"connCount,omitempty"`
}

// WSParseEvent 解析 event 类型消息。
// ok=false 表示该消息不是 event 消息（通常是 data 推送）。
func WSParseEvent(message []byte) (*WSEvent, bool, error) {
	var ev WSEvent
	if err := json.Unmarshal(message, &ev); err != nil {
		return nil, false, err
	}
	if ev.Event == "" {
		return nil, false, nil
	}
	return &ev, true, nil
}

// WSData 表示 OKX WebSocket data 推送。
type WSData[T any] struct {
	Arg       WSArg  `json:"arg"`
	EventType string `json:"eventType,omitempty"`
	Action    string `json:"action,omitempty"`
	CurPage   int    `json:"curPage,omitempty"`
	LastPage  bool   `json:"lastPage,omitempty"`
	Data      []T    `json:"data"`
}

// WSParseData 解析 data 推送消息。
// ok=false 表示该消息不是 data 推送（通常是 event）。
func WSParseData[T any](message []byte) (*WSData[T], bool, error) {
	var dm WSData[T]
	if err := json.Unmarshal(message, &dm); err != nil {
		return nil, false, err
	}
	if dm.Arg.Channel == "" {
		return nil, false, nil
	}
	if dm.Data == nil {
		return nil, false, nil
	}
	return &dm, true, nil
}

// WSLiquidationWarning 表示爆仓风险预警推送（liquidation-warning）。
type WSLiquidationWarning struct {
	InstType string `json:"instType"`
	MgnMode  string `json:"mgnMode"`
	PosId    string `json:"posId"`
	PosSide  string `json:"posSide"`
	Pos      string `json:"pos"`
	PosCcy   string `json:"posCcy"`
	InstId   string `json:"instId"`
	Lever    string `json:"lever"`
	MarkPx   string `json:"markPx"`
	MgnRatio string `json:"mgnRatio"`
	Ccy      string `json:"ccy"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime"`
	PTime UnixMilli `json:"pTime"`
}

// WSADLWarning 表示自动减仓预警推送（adl-warning）。
//
// 说明：时间戳字段可能为空字符串，使用 UnixMilli 兼容解析。
type WSADLWarning struct {
	InstType   string `json:"instType"`
	InstFamily string `json:"instFamily"`
	Ccy        string `json:"ccy"`

	State string `json:"state"`
	Bal   string `json:"bal"`

	MaxBal   string    `json:"maxBal"`
	MaxBalTS UnixMilli `json:"maxBalTs"`

	ADLType string `json:"adlType"`
	ADLBal  string `json:"adlBal"`

	ADLRecBal string `json:"adlRecBal"`

	DecRate    string `json:"decRate"`
	ADLRate    string `json:"adlRate"`
	ADLRecRate string `json:"adlRecRate"`

	TS UnixMilli `json:"ts"`
}

// WSParseChannelData 解析指定 channel 的 data 推送消息。
func WSParseChannelData[T any](message []byte, channel string) (*WSData[T], bool, error) {
	dm, ok, err := WSParseData[T](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if dm.Arg.Channel != channel {
		return nil, false, nil
	}
	return dm, true, nil
}

// WSParseOrders 解析 orders 频道推送消息。
func WSParseOrders(message []byte) (*WSData[TradeOrder], bool, error) {
	return WSParseChannelData[TradeOrder](message, WSChannelOrders)
}

// WSParseFills 解析 fills 频道推送消息。
func WSParseFills(message []byte) (*WSData[WSFill], bool, error) {
	return WSParseChannelData[WSFill](message, WSChannelFills)
}

// WSParseAccount 解析 account 频道推送消息。
func WSParseAccount(message []byte) (*WSData[AccountBalance], bool, error) {
	return WSParseChannelData[AccountBalance](message, WSChannelAccount)
}

// WSParsePositions 解析 positions 频道推送消息。
func WSParsePositions(message []byte) (*WSData[AccountPosition], bool, error) {
	return WSParseChannelData[AccountPosition](message, WSChannelPositions)
}

// WSParseBalanceAndPosition 解析 balance_and_position 频道推送消息。
func WSParseBalanceAndPosition(message []byte) (*WSData[WSBalanceAndPosition], bool, error) {
	return WSParseChannelData[WSBalanceAndPosition](message, WSChannelBalanceAndPosition)
}

// WSParseLiquidationWarning 解析 liquidation-warning 频道推送消息（private WS，需要登录）。
func WSParseLiquidationWarning(message []byte) (*WSData[WSLiquidationWarning], bool, error) {
	return WSParseChannelData[WSLiquidationWarning](message, WSChannelLiquidationWarning)
}

// WSParseAccountGreeks 解析 account-greeks 频道推送消息（private WS，需要登录）。
func WSParseAccountGreeks(message []byte) (*WSData[AccountGreeks], bool, error) {
	return WSParseChannelData[AccountGreeks](message, WSChannelAccountGreeks)
}

// WSParseOrdersAlgo 解析 orders-algo 频道推送消息（business WS，需要登录）。
func WSParseOrdersAlgo(message []byte) (*WSData[TradeAlgoOrder], bool, error) {
	return WSParseChannelData[TradeAlgoOrder](message, WSChannelOrdersAlgo)
}

// WSParseAlgoAdvance 解析 algo-advance 频道推送消息（business WS，需要登录）。
func WSParseAlgoAdvance(message []byte) (*WSData[TradeAlgoOrder], bool, error) {
	return WSParseChannelData[TradeAlgoOrder](message, WSChannelAlgoAdvance)
}

// WSParseGridOrdersSpot 解析 grid-orders-spot 频道推送消息（business WS，需要登录）。
func WSParseGridOrdersSpot(message []byte) (*WSData[WSGridOrder], bool, error) {
	return WSParseChannelData[WSGridOrder](message, WSChannelGridOrdersSpot)
}

// WSParseGridOrdersContract 解析 grid-orders-contract 频道推送消息（business WS，需要登录）。
func WSParseGridOrdersContract(message []byte) (*WSData[WSGridOrder], bool, error) {
	return WSParseChannelData[WSGridOrder](message, WSChannelGridOrdersContract)
}

// WSParseGridPositions 解析 grid-positions 频道推送消息（business WS，需要登录）。
func WSParseGridPositions(message []byte) (*WSData[WSGridPosition], bool, error) {
	return WSParseChannelData[WSGridPosition](message, WSChannelGridPositions)
}

// WSParseGridSubOrders 解析 grid-sub-orders 频道推送消息（business WS，需要登录）。
func WSParseGridSubOrders(message []byte) (*WSData[WSGridSubOrder], bool, error) {
	return WSParseChannelData[WSGridSubOrder](message, WSChannelGridSubOrders)
}

// WSParseAlgoRecurringBuy 解析 algo-recurring-buy 频道推送消息（business WS，需要登录）。
func WSParseAlgoRecurringBuy(message []byte) (*WSData[WSRecurringBuyOrder], bool, error) {
	return WSParseChannelData[WSRecurringBuyOrder](message, WSChannelAlgoRecurringBuy)
}

// WSParseCopytradingLeadNotification 解析 copytrading-lead-notification 频道推送消息（business WS，需要登录）。
func WSParseCopytradingLeadNotification(message []byte) (*WSData[WSCopyTradingLeadNotification], bool, error) {
	return WSParseChannelData[WSCopyTradingLeadNotification](message, WSChannelCopytradingLeadNotification)
}

// WSParseRFQs 解析 rfqs 频道推送消息（business WS，需要登录）。
func WSParseRFQs(message []byte) (*WSData[WSRFQ], bool, error) {
	return WSParseChannelData[WSRFQ](message, WSChannelRFQs)
}

// WSParseQuotes 解析 quotes 频道推送消息（business WS，需要登录）。
func WSParseQuotes(message []byte) (*WSData[WSQuote], bool, error) {
	return WSParseChannelData[WSQuote](message, WSChannelQuotes)
}

// WSParseStrucBlockTrades 解析 struc-block-trades 频道推送消息（business WS，需要登录）。
func WSParseStrucBlockTrades(message []byte) (*WSData[WSStrucBlockTrade], bool, error) {
	return WSParseChannelData[WSStrucBlockTrade](message, WSChannelStrucBlockTrades)
}

// WSParsePublicStrucBlockTrades 解析 public-struc-block-trades 频道推送消息（business WS，无需登录）。
func WSParsePublicStrucBlockTrades(message []byte) (*WSData[WSPublicStrucBlockTrade], bool, error) {
	return WSParseChannelData[WSPublicStrucBlockTrade](message, WSChannelPublicStrucBlockTrades)
}

// WSParsePublicBlockTrades 解析 public-block-trades 频道推送消息（business WS，无需登录）。
func WSParsePublicBlockTrades(message []byte) (*WSData[BlockTrade], bool, error) {
	return WSParseChannelData[BlockTrade](message, WSChannelPublicBlockTrades)
}

// WSParseBlockTickers 解析 block-tickers 频道推送消息（business WS，无需登录）。
func WSParseBlockTickers(message []byte) (*WSData[WSBlockTicker], bool, error) {
	return WSParseChannelData[WSBlockTicker](message, WSChannelBlockTickers)
}

// WSParseDepositInfo 解析 deposit-info 频道推送消息（business WS，需要登录）。
func WSParseDepositInfo(message []byte) (*WSData[WSDepositInfo], bool, error) {
	return WSParseChannelData[WSDepositInfo](message, WSChannelDepositInfo)
}

// WSParseWithdrawalInfo 解析 withdrawal-info 频道推送消息（business WS，需要登录）。
func WSParseWithdrawalInfo(message []byte) (*WSData[WSWithdrawalInfo], bool, error) {
	return WSParseChannelData[WSWithdrawalInfo](message, WSChannelWithdrawalInfo)
}

// WSParseSprdTickers 解析 sprd-tickers 频道推送消息（business WS，无需登录）。
func WSParseSprdTickers(message []byte) (*WSData[MarketSprdTicker], bool, error) {
	return WSParseChannelData[MarketSprdTicker](message, WSChannelSprdTickers)
}

// WSParseSprdOrders 解析 sprd-orders 频道推送消息（business WS，需要登录）。
func WSParseSprdOrders(message []byte) (*WSData[SprdOrder], bool, error) {
	return WSParseChannelData[SprdOrder](message, WSChannelSprdOrders)
}

// WSParseSprdTrades 解析 sprd-trades 频道推送消息（business WS，需要登录）。
func WSParseSprdTrades(message []byte) (*WSData[SprdTrade], bool, error) {
	return WSParseChannelData[SprdTrade](message, WSChannelSprdTrades)
}

// WSParseTickers 解析 tickers 频道推送消息。
func WSParseTickers(message []byte) (*WSData[MarketTicker], bool, error) {
	return WSParseChannelData[MarketTicker](message, WSChannelTickers)
}

// WSParseTrades 解析 trades 频道推送消息。
func WSParseTrades(message []byte) (*WSData[MarketTrade], bool, error) {
	return WSParseChannelData[MarketTrade](message, WSChannelTrades)
}

// WSParseTradesAll 解析 trades-all 频道推送消息。
func WSParseTradesAll(message []byte) (*WSData[MarketTrade], bool, error) {
	return WSParseChannelData[MarketTrade](message, WSChannelTradesAll)
}

// WSParseStatus 解析 status 频道推送消息（public WS）。
func WSParseStatus(message []byte) (*WSData[SystemStatus], bool, error) {
	return WSParseChannelData[SystemStatus](message, WSChannelStatus)
}

// WSParseOpenInterest 解析 open-interest 频道推送消息。
func WSParseOpenInterest(message []byte) (*WSData[OpenInterest], bool, error) {
	return WSParseChannelData[OpenInterest](message, WSChannelOpenInterest)
}

// WSParseFundingRate 解析 funding-rate 频道推送消息。
func WSParseFundingRate(message []byte) (*WSData[FundingRate], bool, error) {
	return WSParseChannelData[FundingRate](message, WSChannelFundingRate)
}

// WSParsePriceLimit 解析 price-limit 频道推送消息。
func WSParsePriceLimit(message []byte) (*WSData[PriceLimit], bool, error) {
	return WSParseChannelData[PriceLimit](message, WSChannelPriceLimit)
}

// WSParseMarkPrice 解析 mark-price 频道推送消息。
func WSParseMarkPrice(message []byte) (*WSData[MarkPrice], bool, error) {
	return WSParseChannelData[MarkPrice](message, WSChannelMarkPrice)
}

// WSParseIndexTickers 解析 index-tickers 频道推送消息。
func WSParseIndexTickers(message []byte) (*WSData[IndexTicker], bool, error) {
	return WSParseChannelData[IndexTicker](message, WSChannelIndexTickers)
}

// WSParseOptSummary 解析 opt-summary 频道推送消息。
func WSParseOptSummary(message []byte) (*WSData[OptSummary], bool, error) {
	return WSParseChannelData[OptSummary](message, WSChannelOptSummary)
}

// WSParseInstruments 解析 instruments 频道推送消息。
func WSParseInstruments(message []byte) (*WSData[Instrument], bool, error) {
	return WSParseChannelData[Instrument](message, WSChannelInstruments)
}

// WSParseEstimatedPrice 解析 estimated-price 频道推送消息。
func WSParseEstimatedPrice(message []byte) (*WSData[EstimatedPrice], bool, error) {
	return WSParseChannelData[EstimatedPrice](message, WSChannelEstimatedPrice)
}

// WSParseADLWarning 解析 adl-warning 频道推送消息。
func WSParseADLWarning(message []byte) (*WSData[WSADLWarning], bool, error) {
	return WSParseChannelData[WSADLWarning](message, WSChannelADLWarning)
}

// WSParseEconomicCalendar 解析 economic-calendar 频道推送消息（business WS，需要登录）。
func WSParseEconomicCalendar(message []byte) (*WSData[EconomicCalendarEvent], bool, error) {
	return WSParseChannelData[EconomicCalendarEvent](message, WSChannelEconomicCalendar)
}

// WSParseLiquidationOrders 解析 liquidation-orders 频道推送消息。
func WSParseLiquidationOrders(message []byte) (*WSData[LiquidationOrder], bool, error) {
	return WSParseChannelData[LiquidationOrder](message, WSChannelLiquidationOrders)
}

// WSCandleChannel 返回 OKX K线频道名（如 bar=1m -> candle1m）。
func WSCandleChannel(bar string) string {
	if bar == "" {
		return ""
	}
	if strings.HasPrefix(bar, "candle") {
		return bar
	}
	return "candle" + bar
}

func isCandleChannel(channel string) bool {
	return strings.HasPrefix(channel, "candle")
}

const wsChannelPrefixSprdCandle = "sprd-candle"

// WSSprdCandleChannel 返回价差交易 K线频道名（如 bar=1D -> sprd-candle1D）。
func WSSprdCandleChannel(bar string) string {
	if bar == "" {
		return ""
	}
	if strings.HasPrefix(bar, wsChannelPrefixSprdCandle) {
		return bar
	}
	return wsChannelPrefixSprdCandle + bar
}

func isSprdCandleChannel(channel string) bool {
	return strings.HasPrefix(channel, wsChannelPrefixSprdCandle)
}

// WSParseCandles 解析 K线频道推送消息（candle*，business WS）。
func WSParseCandles(message []byte) (*WSData[Candle], bool, error) {
	dm, ok, err := WSParseData[Candle](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if !isCandleChannel(dm.Arg.Channel) {
		return nil, false, nil
	}
	return dm, true, nil
}

// WSParseSprdCandles 解析价差交易 K线频道推送消息（sprd-candle*，business WS，无需登录）。
func WSParseSprdCandles(message []byte) (*WSData[Candle], bool, error) {
	dm, ok, err := WSParseData[Candle](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if !isSprdCandleChannel(dm.Arg.Channel) {
		return nil, false, nil
	}
	return dm, true, nil
}

const (
	wsChannelPrefixMarkPriceCandle = "mark-price-candle"
	wsChannelPrefixIndexCandle     = "index-candle"
)

// WSMarkPriceCandleChannel 返回标记价格K线频道名（如 bar=1D -> mark-price-candle1D）。
func WSMarkPriceCandleChannel(bar string) string {
	if bar == "" {
		return ""
	}
	if strings.HasPrefix(bar, wsChannelPrefixMarkPriceCandle) {
		return bar
	}
	return wsChannelPrefixMarkPriceCandle + bar
}

// WSIndexCandleChannel 返回指数K线频道名（如 bar=30m -> index-candle30m）。
func WSIndexCandleChannel(bar string) string {
	if bar == "" {
		return ""
	}
	if strings.HasPrefix(bar, wsChannelPrefixIndexCandle) {
		return bar
	}
	return wsChannelPrefixIndexCandle + bar
}

func isMarkPriceCandleChannel(channel string) bool {
	return strings.HasPrefix(channel, wsChannelPrefixMarkPriceCandle)
}

func isIndexCandleChannel(channel string) bool {
	return strings.HasPrefix(channel, wsChannelPrefixIndexCandle)
}

// WSParseMarkPriceCandles 解析标记价格K线频道推送消息（mark-price-candle*，business WS）。
func WSParseMarkPriceCandles(message []byte) (*WSData[PriceCandle], bool, error) {
	dm, ok, err := WSParseData[PriceCandle](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if !isMarkPriceCandleChannel(dm.Arg.Channel) {
		return nil, false, nil
	}
	return dm, true, nil
}

// WSParseIndexCandles 解析指数K线频道推送消息（index-candle*，business WS）。
func WSParseIndexCandles(message []byte) (*WSData[PriceCandle], bool, error) {
	dm, ok, err := WSParseData[PriceCandle](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if !isIndexCandleChannel(dm.Arg.Channel) {
		return nil, false, nil
	}
	return dm, true, nil
}

// WSCandle 表示 K线推送的单条数据项（含订阅 Arg 上下文）。
//
// 说明：OKX K线数据本身不包含 instId/sprdId，因此需要通过 Arg 携带产品信息。
type WSCandle struct {
	Arg    WSArg
	Candle Candle
}

// WSPriceCandle 表示指数/标记价格 K线推送的单条数据项（含订阅 Arg 上下文）。
//
// 说明：OKX PriceCandle 数据本身不包含 instId，因此需要通过 Arg 携带产品信息。
type WSPriceCandle struct {
	Arg    WSArg
	Candle PriceCandle
}

// WSOptionTrade 表示 option-trades 频道推送的数据项。
type WSOptionTrade struct {
	InstId     string `json:"instId"`
	InstFamily string `json:"instFamily"`

	TradeId string `json:"tradeId"`

	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Side string `json:"side"`

	OptType string `json:"optType"`
	FillVol string `json:"fillVol"`
	FwdPx   string `json:"fwdPx"`
	IdxPx   string `json:"idxPx"`
	MarkPx  string `json:"markPx"`

	TS int64 `json:"ts,string"`
}

// WSParseOptionTrades 解析 option-trades 频道推送消息。
func WSParseOptionTrades(message []byte) (*WSData[WSOptionTrade], bool, error) {
	return WSParseChannelData[WSOptionTrade](message, WSChannelOptionTrades)
}

// WSCallAuctionDetails 表示 call-auction-details 频道推送的数据项。
type WSCallAuctionDetails struct {
	InstId string `json:"instId"`

	EqPx        string `json:"eqPx"`
	MatchedSz   string `json:"matchedSz"`
	UnmatchedSz string `json:"unmatchedSz"`

	State          string `json:"state"`
	AuctionEndTime int64  `json:"auctionEndTime,string"`
	TS             int64  `json:"ts,string"`
}

// WSParseCallAuctionDetails 解析 call-auction-details 频道推送消息。
func WSParseCallAuctionDetails(message []byte) (*WSData[WSCallAuctionDetails], bool, error) {
	return WSParseChannelData[WSCallAuctionDetails](message, WSChannelCallAuctionDetails)
}

// WSSprdPublicTrade 表示 sprd-public-trades 频道推送的数据项。
//
// 说明：价格/数量等字段保持为 string（无损）。
type WSSprdPublicTrade struct {
	SprdId  string `json:"sprdId"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	Px      string `json:"px"`
	TradeId string `json:"tradeId"`

	TS int64 `json:"ts,string"`
}

// WSParseSprdPublicTrades 解析 sprd-public-trades 频道推送消息（business WS）。
func WSParseSprdPublicTrades(message []byte) (*WSData[WSSprdPublicTrade], bool, error) {
	return WSParseChannelData[WSSprdPublicTrade](message, WSChannelSprdPublicTrades)
}

// PriceLimit 表示 WS price-limit 频道推送的数据项。
type PriceLimit struct {
	InstType string `json:"instType,omitempty"`
	InstId   string `json:"instId"`

	BuyLmt  string `json:"buyLmt"`
	SellLmt string `json:"sellLmt"`

	TS int64 `json:"ts,string"`

	Enabled bool `json:"enabled"`
}

// IndexTicker 表示 WS index-tickers 频道推送的数据项。
type IndexTicker struct {
	InstId string `json:"instId"`
	IdxPx  string `json:"idxPx"`

	High24h string `json:"high24h"`
	Low24h  string `json:"low24h"`
	Open24h string `json:"open24h"`

	SodUtc0 string `json:"sodUtc0"`
	SodUtc8 string `json:"sodUtc8"`

	TS int64 `json:"ts,string"`
}

// LiquidationOrder 表示 WS liquidation-orders 频道推送的数据项。
type LiquidationOrder struct {
	InstType   string `json:"instType"`
	InstId     string `json:"instId"`
	Uly        string `json:"uly,omitempty"`
	InstFamily string `json:"instFamily,omitempty"`

	Details []LiquidationOrderDetail `json:"details"`
}

type LiquidationOrderDetail struct {
	Side    string `json:"side"`
	PosSide string `json:"posSide"`

	BkPx   string `json:"bkPx"`
	Sz     string `json:"sz"`
	BkLoss string `json:"bkLoss"`
	Ccy    string `json:"ccy"`

	TS int64 `json:"ts,string"`
}

// WSOrderBook 表示 WS 深度频道推送的数据项。
type WSOrderBook struct {
	Asks []OrderBookLevel `json:"asks"`
	Bids []OrderBookLevel `json:"bids"`

	InstId string `json:"instId"`

	TS int64 `json:"ts,string"`

	Checksum  int64 `json:"checksum,omitempty"`
	PrevSeqId int64 `json:"prevSeqId,omitempty"`
	SeqId     int64 `json:"seqId,omitempty"`
}

func isOrderBookChannel(channel string) bool {
	switch channel {
	case WSChannelBooks, WSChannelBooksELP, WSChannelBooks5, WSChannelBboTbt, WSChannelBooksL2Tbt, WSChannelBooks50L2Tbt,
		WSChannelSprdBboTbt, WSChannelSprdBooks5, WSChannelSprdBooksL2Tbt:
		return true
	default:
		return false
	}
}

// WSParseOrderBook 解析深度频道推送消息（books/books5/bbo-tbt/books-l2-tbt/books50-l2-tbt/books-elp）。
func WSParseOrderBook(message []byte) (*WSData[WSOrderBook], bool, error) {
	dm, ok, err := WSParseData[WSOrderBook](message)
	if err != nil || !ok {
		return nil, ok, err
	}
	if !isOrderBookChannel(dm.Arg.Channel) {
		return nil, false, nil
	}
	return dm, true, nil
}

// WSBalanceAndPosition 表示 balance_and_position 频道推送的数据项（精简版）。
type WSBalanceAndPosition struct {
	PTime     int64  `json:"pTime,string"`
	EventType string `json:"eventType"`

	BalData []WSBalanceAndPositionBalance `json:"balData"`
	PosData []AccountPosition             `json:"posData"`
}

type WSBalanceAndPositionBalance struct {
	Ccy     string `json:"ccy"`
	CashBal string `json:"cashBal"`
	UTime   int64  `json:"uTime,string"`
}

// WSFill 表示 WS / 成交频道推送的数据项。
// 该频道仅适用于交易等级 VIP6 及以上用户；其他用户可使用 orders 频道获取成交信息。
type WSFill struct {
	InstId string `json:"instId"`
	FillSz string `json:"fillSz"`
	FillPx string `json:"fillPx"`
	Side   string `json:"side"`

	TS string `json:"ts"`

	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	TradeId string `json:"tradeId"`

	ExecType string `json:"execType"`
	Count    string `json:"count"`
}

// WSDepositInfo 表示充值信息推送（deposit-info）。
type WSDepositInfo struct {
	AssetDeposit

	PTime   int64  `json:"pTime,string"`
	SubAcct string `json:"subAcct"`
	UID     string `json:"uid"`
}

// WSWithdrawalInfo 表示提币信息推送（withdrawal-info）。
type WSWithdrawalInfo struct {
	AssetWithdrawal

	PTime   int64  `json:"pTime,string"`
	SubAcct string `json:"subAcct"`
	UID     string `json:"uid"`
}

// WSGridOrder 表示网格策略委托订单推送（grid-orders-spot / grid-orders-contract）的数据项（精简版）。
//
// 说明：数值字段保持为 string（无损），时间戳字段解析为 int64。
type WSGridOrder struct {
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoId      string `json:"algoId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	State string `json:"state"`
	Tag   string `json:"tag"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`
	PTime int64 `json:"pTime,string"`

	GridNum      string `json:"gridNum"`
	ActiveOrdNum string `json:"activeOrdNum"`

	MinPx string `json:"minPx"`
	MaxPx string `json:"maxPx"`
	RunPx string `json:"runPx"`

	RunType    string `json:"runType"`
	CancelType string `json:"cancelType"`
	StopType   string `json:"stopType"`

	Investment string `json:"investment"`
	Sz         string `json:"sz"`

	Profit      string `json:"profit"`
	FloatProfit string `json:"floatProfit"`
	GridProfit  string `json:"gridProfit"`
	TotalPnl    string `json:"totalPnl"`
	PnlRatio    string `json:"pnlRatio"`

	AnnualizedRate      string `json:"annualizedRate"`
	TotalAnnualizedRate string `json:"totalAnnualizedRate"`

	PerMaxProfitRate string `json:"perMaxProfitRate"`
	PerMinProfitRate string `json:"perMinProfitRate"`

	Lever       string `json:"lever"`
	ActualLever string `json:"actualLever"`
	Direction   string `json:"direction"`
	BasePos     bool   `json:"basePos"`
	AvailEq     string `json:"availEq"`
	Eq          string `json:"eq"`
	LiqPx       string `json:"liqPx"`

	SingleAmt   string `json:"singleAmt"`
	SlTriggerPx string `json:"slTriggerPx"`
	TpTriggerPx string `json:"tpTriggerPx"`

	TradeNum string `json:"tradeNum"`
}

// WSGridPosition 表示网格策略持仓推送（grid-positions）的数据项（精简版）。
//
// 说明：该频道 subscribe 参数需要提供 arg.algoId。
type WSGridPosition struct {
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoId      string `json:"algoId"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	Ccy string `json:"ccy"`

	MgnMode  string `json:"mgnMode"`
	PosSide  string `json:"posSide"`
	Pos      string `json:"pos"`
	Lever    string `json:"lever"`
	AvgPx    string `json:"avgPx"`
	Last     string `json:"last"`
	MarkPx   string `json:"markPx"`
	LiqPx    string `json:"liqPx"`
	MgnRatio string `json:"mgnRatio"`

	IMR         string `json:"imr"`
	MMR         string `json:"mmr"`
	NotionalUSD string `json:"notionalUsd"`

	ADL      string `json:"adl"`
	Upl      string `json:"upl"`
	UplRatio string `json:"uplRatio"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`
	PTime int64 `json:"pTime,string"`
}

// WSCopyTradingLeadNotification 表示带单失败通知推送（copytrading-lead-notification）的数据项（精简版）。
type WSCopyTradingLeadNotification struct {
	InfoType string `json:"infoType"`

	SubPosId         string `json:"subPosId"`
	UniqueCode       string `json:"uniqueCode"`
	InstType         string `json:"instType"`
	InstId           string `json:"instId"`
	Side             string `json:"side"`
	PosSide          string `json:"posSide"`
	MaxLeadTraderNum string `json:"maxLeadTraderNum"`
	MinLeadEq        string `json:"minLeadEq"`
}

// WSGridSubOrder 表示网格策略子订单推送（grid-sub-orders）的数据项（精简版）。
//
// 说明：该频道 subscribe 参数需要提供 arg.algoId。
type WSGridSubOrder struct {
	AlgoId      string `json:"algoId"`
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoOrdType string `json:"algoOrdType"`

	InstType string `json:"instType"`
	InstId   string `json:"instId"`

	OrdId   string `json:"ordId"`
	OrdType string `json:"ordType"`

	Side    string `json:"side"`
	PosSide string `json:"posSide"`
	TdMode  string `json:"tdMode"`

	State string `json:"state"`

	Sz        string `json:"sz"`
	Px        string `json:"px"`
	AvgPx     string `json:"avgPx"`
	AccFillSz string `json:"accFillSz"`

	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	Rebate    string `json:"rebate"`
	RebateCcy string `json:"rebateCcy"`

	Lever   string `json:"lever"`
	GroupId string `json:"groupId"`
	CtVal   string `json:"ctVal"`
	Pnl     string `json:"pnl"`

	Tag string `json:"tag"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`
	PTime int64 `json:"pTime,string"`
}

// WSRecurringBuyOrder 表示定投策略委托订单推送（algo-recurring-buy）的数据项（精简版）。
type WSRecurringBuyOrder struct {
	AlgoClOrdId string `json:"algoClOrdId"`
	AlgoId      string `json:"algoId"`
	AlgoOrdType string `json:"algoOrdType"`
	Tag         string `json:"tag"`

	InstType string `json:"instType"`

	Amt           string `json:"amt"`
	Cycles        string `json:"cycles"`
	InvestmentAmt string `json:"investmentAmt"`
	InvestmentCcy string `json:"investmentCcy"`
	MktCap        string `json:"mktCap"`

	NextInvestTime int64 `json:"nextInvestTime,string"`
	PTime          int64 `json:"pTime,string"`

	Period        string `json:"period"`
	PnlRatio      string `json:"pnlRatio"`
	RecurringDay  string `json:"recurringDay"`
	RecurringHour string `json:"recurringHour"`
	RecurringTime string `json:"recurringTime"`
	TimeZone      string `json:"timeZone"`

	State    string `json:"state"`
	StgyName string `json:"stgyName"`

	TotalAnnRate string `json:"totalAnnRate"`
	TotalPnl     string `json:"totalPnl"`

	TradeQuoteCcy string `json:"tradeQuoteCcy"`

	RecurringList []WSRecurringBuyItem `json:"recurringList"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`
}

type WSRecurringBuyItem struct {
	Ccy      string `json:"ccy"`
	Ratio    string `json:"ratio"`
	Px       string `json:"px"`
	AvgPx    string `json:"avgPx"`
	Profit   string `json:"profit"`
	TotalAmt string `json:"totalAmt"`
}

// WSRFQ 表示询价单推送（rfqs）的数据项（精简版）。
type WSRFQ struct {
	RfqId   string `json:"rfqId"`
	ClRfqId string `json:"clRfqId"`
	Tag     string `json:"tag"`

	CTime int64 `json:"cTime,string"`
	UTime int64 `json:"uTime,string"`

	State      string `json:"state"`
	FlowType   string `json:"flowType"`
	TraderCode string `json:"traderCode"`

	ValidUntil            int64    `json:"validUntil,string"`
	AllowPartialExecution bool     `json:"allowPartialExecution"`
	Counterparties        []string `json:"counterparties"`

	Legs []WSRFQLeg `json:"legs"`

	GroupId   string           `json:"groupId"`
	AcctAlloc []WSRFQAcctAlloc `json:"acctAlloc"`
}

type WSRFQLeg struct {
	InstId  string `json:"instId"`
	TdMode  string `json:"tdMode,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
	Sz      string `json:"sz"`
	Side    string `json:"side"`
	PosSide string `json:"posSide,omitempty"`
	TgtCcy  string `json:"tgtCcy,omitempty"`

	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}

type WSRFQAcctAlloc struct {
	Acct  string              `json:"acct"`
	SCode string              `json:"sCode,omitempty"`
	SMsg  string              `json:"sMsg,omitempty"`
	Legs  []WSRFQAcctAllocLeg `json:"legs"`
}

type WSRFQAcctAllocLeg struct {
	InstId  string `json:"instId"`
	Sz      string `json:"sz"`
	TdMode  string `json:"tdMode,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
	PosSide string `json:"posSide,omitempty"`
}

// WSQuote 表示报价单推送（quotes）的数据项（精简版）。
type WSQuote struct {
	QuoteId   string `json:"quoteId"`
	ClQuoteId string `json:"clQuoteId"`
	RfqId     string `json:"rfqId"`
	ClRfqId   string `json:"clRfqId"`
	Tag       string `json:"tag"`

	ValidUntil int64 `json:"validUntil,string"`
	CTime      int64 `json:"cTime,string"`
	UTime      int64 `json:"uTime,string"`

	TraderCode string `json:"traderCode"`
	QuoteSide  string `json:"quoteSide"`
	State      string `json:"state"`
	Reason     string `json:"reason"`

	Legs []WSQuoteLeg `json:"legs"`
}

type WSQuoteLeg struct {
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	InstId string `json:"instId"`

	TdMode  string `json:"tdMode,omitempty"`
	Ccy     string `json:"ccy,omitempty"`
	Side    string `json:"side"`
	PosSide string `json:"posSide,omitempty"`
	TgtCcy  string `json:"tgtCcy,omitempty"`

	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}

// WSStrucBlockTrade 表示大宗交易推送（struc-block-trades）的数据项（精简版）。
type WSStrucBlockTrade struct {
	CTime int64 `json:"cTime,string"`

	RfqId     string `json:"rfqId"`
	ClRfqId   string `json:"clRfqId"`
	QuoteId   string `json:"quoteId"`
	ClQuoteId string `json:"clQuoteId"`
	BlockTdId string `json:"blockTdId"`
	GroupId   string `json:"groupId"`
	Tag       string `json:"tag"`

	TTraderCode string `json:"tTraderCode"`
	MTraderCode string `json:"mTraderCode"`

	IsSuccessful bool   `json:"isSuccessful"`
	ErrorCode    string `json:"errorCode"`

	Legs      []WSStrucBlockTradeLeg       `json:"legs"`
	AcctAlloc []WSStrucBlockTradeAcctAlloc `json:"acctAlloc,omitempty"`
}

// WSPublicStrucBlockTrade 表示公共大宗结构化成交推送（public-struc-block-trades）的数据项。
type WSPublicStrucBlockTrade struct {
	CTime int64 `json:"cTime,string"`

	BlockTdId string `json:"blockTdId"`
	GroupId   string `json:"groupId"`

	Legs []WSPublicStrucBlockTradeLeg `json:"legs"`
}

type WSPublicStrucBlockTradeLeg struct {
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	InstId string `json:"instId"`
	Side   string `json:"side"`

	TradeId string `json:"tradeId"`
}

// WSBlockTicker 表示大宗交易行情推送（block-tickers）的数据项。
type WSBlockTicker = MarketBlockTicker

type WSStrucBlockTradeLeg struct {
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	InstId string `json:"instId"`
	Side   string `json:"side"`

	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	TradeId       string `json:"tradeId"`
	TgtCcy        string `json:"tgtCcy,omitempty"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}

type WSStrucBlockTradeAcctAlloc struct {
	BlockTdId string                          `json:"blockTdId"`
	ErrorCode string                          `json:"errorCode"`
	Acct      string                          `json:"acct"`
	Legs      []WSStrucBlockTradeAcctAllocLeg `json:"legs"`
}

type WSStrucBlockTradeAcctAllocLeg struct {
	InstId string `json:"instId"`
	Px     string `json:"px"`
	Sz     string `json:"sz"`
	Side   string `json:"side"`

	Fee    string `json:"fee"`
	FeeCcy string `json:"feeCcy"`

	TradeId       string `json:"tradeId"`
	TradeQuoteCcy string `json:"tradeQuoteCcy,omitempty"`
}
