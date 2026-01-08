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

	WSChannelTickers   = "tickers"
	WSChannelTrades    = "trades"
	WSChannelTradesAll = "trades-all"

	WSChannelBooks        = "books"
	WSChannelBooksELP     = "books-elp"
	WSChannelBooks5       = "books5"
	WSChannelBboTbt       = "bbo-tbt"
	WSChannelBooksL2Tbt   = "books-l2-tbt"
	WSChannelBooks50L2Tbt = "books50-l2-tbt"

	WSChannelOptionTrades       = "option-trades"
	WSChannelCallAuctionDetails = "call-auction-details"
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
	case WSChannelBooks, WSChannelBooksELP, WSChannelBooks5, WSChannelBboTbt, WSChannelBooksL2Tbt, WSChannelBooks50L2Tbt:
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
