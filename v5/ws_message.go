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

	WSChannelDepositInfo    = "deposit-info"
	WSChannelWithdrawalInfo = "withdrawal-info"

	WSChannelTickers   = "tickers"
	WSChannelTrades    = "trades"
	WSChannelTradesAll = "trades-all"

	WSChannelOpenInterest = "open-interest"
	WSChannelFundingRate  = "funding-rate"
	WSChannelPriceLimit   = "price-limit"
	WSChannelMarkPrice    = "mark-price"
	WSChannelIndexTickers = "index-tickers"
	WSChannelOptSummary   = "opt-summary"

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

// WSParseDepositInfo 解析 deposit-info 频道推送消息（business WS，需要登录）。
func WSParseDepositInfo(message []byte) (*WSData[WSDepositInfo], bool, error) {
	return WSParseChannelData[WSDepositInfo](message, WSChannelDepositInfo)
}

// WSParseWithdrawalInfo 解析 withdrawal-info 频道推送消息（business WS，需要登录）。
func WSParseWithdrawalInfo(message []byte) (*WSData[WSWithdrawalInfo], bool, error) {
	return WSParseChannelData[WSWithdrawalInfo](message, WSChannelWithdrawalInfo)
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
