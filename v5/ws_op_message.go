package okx

import "encoding/json"

// WSOpReply 表示 OKX WebSocket 的业务操作响应（如 order/cancel-order/amend-order 等）。
//
// 这类消息不同于 event（subscribe/login/error/notice），也不同于 data 推送（arg+data）。
type WSOpReply struct {
	ID string `json:"id,omitempty"`
	Op string `json:"op,omitempty"`

	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`

	Data json.RawMessage `json:"data,omitempty"`

	InTime  string `json:"inTime,omitempty"`
	OutTime string `json:"outTime,omitempty"`
}

// WSParseOpReply 解析 op 类型响应消息。
// ok=false 表示该消息不是 op 响应（通常是 event 或 data 推送）。
func WSParseOpReply(message []byte) (*WSOpReply, bool, error) {
	var r WSOpReply
	if err := json.Unmarshal(message, &r); err != nil {
		return nil, false, err
	}
	if r.ID == "" || r.Op == "" {
		return nil, false, nil
	}
	return &r, true, nil
}
