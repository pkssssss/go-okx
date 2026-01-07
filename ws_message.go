package okx

import (
	"encoding/json"
)

const (
	WSChannelOrders = "orders"
	WSChannelFills  = "fills"
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
	Arg  WSArg `json:"arg"`
	Data []T   `json:"data"`
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
