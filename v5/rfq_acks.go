package okx

// RFQTsAck 表示仅返回 ts 的通用响应项（RFQ 模块）。
type RFQTsAck struct {
	TS int64 `json:"ts,string"`
}
