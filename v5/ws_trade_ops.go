package okx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	wsOpOrder       = "order"
	wsOpCancelOrder = "cancel-order"
	wsOpAmendOrder  = "amend-order"
)

var errWSPrivateRequired = errors.New("okx: ws private client required")

// WSPlaceOrderArg 表示 WS / 下单 的 args 项（精简版）。
// 数值字段保持为 string（无损）。
type WSPlaceOrderArg struct {
	InstIdCode int64  `json:"instIdCode,omitempty"`
	InstId     string `json:"instId,omitempty"`

	TdMode string `json:"tdMode"`

	Ccy     string `json:"ccy,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`

	Side    string `json:"side"`
	PosSide string `json:"posSide,omitempty"`

	OrdType string `json:"ordType"`

	Px     string `json:"px,omitempty"`
	PxUsd  string `json:"pxUsd,omitempty"`
	PxVol  string `json:"pxVol,omitempty"`
	PxType string `json:"pxType,omitempty"`

	Sz string `json:"sz"`

	ReduceOnly *bool  `json:"reduceOnly,omitempty"`
	TgtCcy     string `json:"tgtCcy,omitempty"`
	ExpTime    string `json:"expTime,omitempty"`
}

// WSCancelOrderArg 表示 WS / 撤单 的 args 项。
type WSCancelOrderArg struct {
	InstIdCode int64  `json:"instIdCode,omitempty"`
	InstId     string `json:"instId,omitempty"`

	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

// WSAmendOrderArg 表示 WS / 改单 的 args 项（精简版）。
type WSAmendOrderArg struct {
	InstIdCode int64  `json:"instIdCode,omitempty"`
	InstId     string `json:"instId,omitempty"`

	CxlOnFail *bool  `json:"cxlOnFail,omitempty"`
	OrdId     string `json:"ordId,omitempty"`
	ClOrdId   string `json:"clOrdId,omitempty"`

	ReqId string `json:"reqId,omitempty"`

	NewSz    string `json:"newSz,omitempty"`
	NewPx    string `json:"newPx,omitempty"`
	NewPxUsd string `json:"newPxUsd,omitempty"`
	NewPxVol string `json:"newPxVol,omitempty"`

	NewPxType string `json:"newPxType,omitempty"`
	ExpTime   string `json:"expTime,omitempty"`
}

// WSTradeOpError 表示 WS 交易 op 的错误（顶层 code!=0 或 data 内 sCode!=0）。
type WSTradeOpError struct {
	ID string
	Op string

	Code string
	Msg  string

	SCode string
	SMsg  string

	InTime  string
	OutTime string

	Raw []byte
}

func (e *WSTradeOpError) Error() string {
	if e == nil {
		return "<OKX WSTradeOpError>"
	}
	if e.SCode != "" || e.SMsg != "" {
		return fmt.Sprintf("<OKX WSTradeOpError> op=%s id=%s code=%s msg=%s sCode=%s sMsg=%s", e.Op, e.ID, e.Code, e.Msg, e.SCode, e.SMsg)
	}
	return fmt.Sprintf("<OKX WSTradeOpError> op=%s id=%s code=%s msg=%s", e.Op, e.ID, e.Code, e.Msg)
}

// WSTradeOpBatchError 表示 WS 交易 op 的批量部分失败（顶层 code=0，但 data[i].sCode!=0）。
type WSTradeOpBatchError struct {
	ID string
	Op string

	Code string
	Msg  string

	InTime  string
	OutTime string

	Acks []TradeOrderAck
	Raw  []byte
}

func (e *WSTradeOpBatchError) Error() string {
	if e == nil {
		return "<OKX WSTradeOpBatchError>"
	}

	failed := 0
	firstCode := ""
	firstMsg := ""
	for _, ack := range e.Acks {
		if ack.SCode != "" && ack.SCode != "0" {
			failed++
			if firstCode == "" {
				firstCode = ack.SCode
				firstMsg = ack.SMsg
			}
		}
	}

	if failed == 0 {
		return fmt.Sprintf("<OKX WSTradeOpBatchError> op=%s id=%s code=%s msg=%s", e.Op, e.ID, e.Code, e.Msg)
	}
	return fmt.Sprintf("<OKX WSTradeOpBatchError> op=%s id=%s failed=%d sCode=%s sMsg=%s code=%s msg=%s", e.Op, e.ID, failed, firstCode, firstMsg, e.Code, e.Msg)
}

func (w *WSClient) requirePrivate() error {
	if w == nil {
		return errors.New("okx: nil ws client")
	}
	if w.kind != wsKindPrivate || !w.needLogin {
		return errWSPrivateRequired
	}
	return nil
}

func validateWSPlaceOrderArg(prefix string, arg WSPlaceOrderArg) error {
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return fmt.Errorf("%s requires instId or instIdCode", prefix)
	}
	if arg.TdMode == "" {
		return fmt.Errorf("%s requires tdMode", prefix)
	}
	if arg.Side == "" {
		return fmt.Errorf("%s requires side", prefix)
	}
	if arg.OrdType == "" {
		return fmt.Errorf("%s requires ordType", prefix)
	}
	if arg.Sz == "" {
		return fmt.Errorf("%s requires sz", prefix)
	}
	if requiresPriceForOrderType(arg.OrdType) && arg.Px == "" && arg.PxUsd == "" && arg.PxVol == "" {
		return fmt.Errorf("%s requires px/pxUsd/pxVol for this ordType", prefix)
	}
	return nil
}

func validateWSCancelOrderArg(prefix string, arg WSCancelOrderArg) error {
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return fmt.Errorf("%s requires instId or instIdCode", prefix)
	}
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return fmt.Errorf("%s requires ordId or clOrdId", prefix)
	}
	return nil
}

func validateWSAmendOrderArg(prefix string, arg WSAmendOrderArg) error {
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return fmt.Errorf("%s requires instId or instIdCode", prefix)
	}
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return fmt.Errorf("%s requires ordId or clOrdId", prefix)
	}
	if arg.NewSz == "" && arg.NewPx == "" && arg.NewPxUsd == "" && arg.NewPxVol == "" {
		return fmt.Errorf("%s requires newSz/newPx/newPxUsd/newPxVol", prefix)
	}
	return nil
}

func wsTradeCheckBatchAcks(reply *WSOpReply, raw []byte, acks []TradeOrderAck) error {
	for _, ack := range acks {
		if ack.SCode != "" && ack.SCode != "0" {
			return &WSTradeOpBatchError{
				ID:      reply.ID,
				Op:      reply.Op,
				Code:    reply.Code,
				Msg:     reply.Msg,
				InTime:  reply.InTime,
				OutTime: reply.OutTime,
				Acks:    acks,
				Raw:     raw,
			}
		}
	}
	return nil
}

// PlaceOrder 通过 WS 下单（op=order）。
// 注意：该方法会真实提交订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) PlaceOrder(ctx context.Context, arg WSPlaceOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if err := validateWSPlaceOrderArg("okx: ws place order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpOrder, []WSPlaceOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws place order empty response data")
	}
	if acks[0].SCode != "" && acks[0].SCode != "0" {
		return nil, &WSTradeOpError{
			ID:      reply.ID,
			Op:      reply.Op,
			Code:    reply.Code,
			Msg:     reply.Msg,
			SCode:   acks[0].SCode,
			SMsg:    acks[0].SMsg,
			InTime:  reply.InTime,
			OutTime: reply.OutTime,
			Raw:     raw,
		}
	}
	return &acks[0], nil
}

// PlaceOrders 通过 WS 批量下单（op=order，args 为数组）。
// 注意：该方法会真实提交订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) PlaceOrders(ctx context.Context, args ...WSPlaceOrderArg) ([]TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("okx: ws place orders requires at least one arg")
	}
	if len(args) > tradeBatchMaxOrders {
		return nil, errors.New("okx: ws place orders max 20 orders")
	}
	for i, arg := range args {
		if err := validateWSPlaceOrderArg(fmt.Sprintf("okx: ws place orders[%d]", i), arg); err != nil {
			return nil, err
		}
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpOrder, args)
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws place orders empty response data")
	}
	if err := wsTradeCheckBatchAcks(reply, raw, acks); err != nil {
		return acks, err
	}
	return acks, nil
}

// CancelOrder 通过 WS 撤单（op=cancel-order）。
func (w *WSClient) CancelOrder(ctx context.Context, arg WSCancelOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if err := validateWSCancelOrderArg("okx: ws cancel order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpCancelOrder, []WSCancelOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws cancel order empty response data")
	}
	if acks[0].SCode != "" && acks[0].SCode != "0" {
		return nil, &WSTradeOpError{
			ID:      reply.ID,
			Op:      reply.Op,
			Code:    reply.Code,
			Msg:     reply.Msg,
			SCode:   acks[0].SCode,
			SMsg:    acks[0].SMsg,
			InTime:  reply.InTime,
			OutTime: reply.OutTime,
			Raw:     raw,
		}
	}
	return &acks[0], nil
}

// CancelOrders 通过 WS 批量撤单（op=cancel-order，args 为数组）。
// 注意：该方法会真实撤销订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) CancelOrders(ctx context.Context, args ...WSCancelOrderArg) ([]TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("okx: ws cancel orders requires at least one arg")
	}
	if len(args) > tradeBatchMaxOrders {
		return nil, errors.New("okx: ws cancel orders max 20 orders")
	}
	for i, arg := range args {
		if err := validateWSCancelOrderArg(fmt.Sprintf("okx: ws cancel orders[%d]", i), arg); err != nil {
			return nil, err
		}
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpCancelOrder, args)
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws cancel orders empty response data")
	}
	if err := wsTradeCheckBatchAcks(reply, raw, acks); err != nil {
		return acks, err
	}
	return acks, nil
}

// AmendOrder 通过 WS 改单（op=amend-order）。
func (w *WSClient) AmendOrder(ctx context.Context, arg WSAmendOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if err := validateWSAmendOrderArg("okx: ws amend order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpAmendOrder, []WSAmendOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws amend order empty response data")
	}
	if acks[0].SCode != "" && acks[0].SCode != "0" {
		return nil, &WSTradeOpError{
			ID:      reply.ID,
			Op:      reply.Op,
			Code:    reply.Code,
			Msg:     reply.Msg,
			SCode:   acks[0].SCode,
			SMsg:    acks[0].SMsg,
			InTime:  reply.InTime,
			OutTime: reply.OutTime,
			Raw:     raw,
		}
	}
	return &acks[0], nil
}

// AmendOrders 通过 WS 批量改单（op=amend-order，args 为数组）。
// 注意：该方法会真实修改订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) AmendOrders(ctx context.Context, args ...WSAmendOrderArg) ([]TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("okx: ws amend orders requires at least one arg")
	}
	if len(args) > tradeBatchMaxOrders {
		return nil, errors.New("okx: ws amend orders max 20 orders")
	}
	for i, arg := range args {
		if err := validateWSAmendOrderArg(fmt.Sprintf("okx: ws amend orders[%d]", i), arg); err != nil {
			return nil, err
		}
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpAmendOrder, args)
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws amend orders empty response data")
	}
	if err := wsTradeCheckBatchAcks(reply, raw, acks); err != nil {
		return acks, err
	}
	return acks, nil
}

func unmarshalTradeOrderAcks(reply *WSOpReply, raw []byte) ([]TradeOrderAck, error) {
	if reply == nil {
		return nil, errors.New("okx: nil ws op reply")
	}
	if reply.Code != "" && reply.Code != "0" {
		return nil, &WSTradeOpError{
			ID:      reply.ID,
			Op:      reply.Op,
			Code:    reply.Code,
			Msg:     reply.Msg,
			InTime:  reply.InTime,
			OutTime: reply.OutTime,
			Raw:     raw,
		}
	}

	if len(reply.Data) == 0 || string(reply.Data) == "null" {
		return nil, nil
	}

	var acks []TradeOrderAck
	if err := json.Unmarshal(reply.Data, &acks); err != nil {
		return nil, err
	}
	return acks, nil
}

// requiresPriceForOrderType 是对 OKX ordType 价格字段要求的“best-effort”本地校验。
// 注意：OKX 文档/服务端语义为准；对于未知 ordType，这里默认不拦截，由服务端返回错误。
func requiresPriceForOrderType(ordType string) bool {
	switch ordType {
	case "limit", "post_only", "fok", "ioc", "mmp", "mmp_and_post_only":
		return true
	default:
		return false
	}
}
