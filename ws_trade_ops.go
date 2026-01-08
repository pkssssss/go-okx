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

func (w *WSClient) requirePrivate() error {
	if w == nil {
		return errors.New("okx: nil ws client")
	}
	if !w.needLogin {
		return errWSPrivateRequired
	}
	return nil
}

// PlaceOrder 通过 WS 下单（op=order）。
// 注意：该方法会真实提交订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) PlaceOrder(ctx context.Context, arg WSPlaceOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return nil, errors.New("okx: ws place order requires instId or instIdCode")
	}
	if arg.TdMode == "" {
		return nil, errors.New("okx: ws place order requires tdMode")
	}
	if arg.Side == "" {
		return nil, errors.New("okx: ws place order requires side")
	}
	if arg.OrdType == "" {
		return nil, errors.New("okx: ws place order requires ordType")
	}
	if arg.Sz == "" {
		return nil, errors.New("okx: ws place order requires sz")
	}
	if requiresPriceForOrderType(arg.OrdType) && arg.Px == "" && arg.PxUsd == "" && arg.PxVol == "" {
		return nil, errors.New("okx: ws place order requires px/pxUsd/pxVol for this ordType")
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

// CancelOrder 通过 WS 撤单（op=cancel-order）。
func (w *WSClient) CancelOrder(ctx context.Context, arg WSCancelOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return nil, errors.New("okx: ws cancel order requires instId or instIdCode")
	}
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return nil, errors.New("okx: ws cancel order requires ordId or clOrdId")
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

// AmendOrder 通过 WS 改单（op=amend-order）。
func (w *WSClient) AmendOrder(ctx context.Context, arg WSAmendOrderArg) (*TradeOrderAck, error) {
	if err := w.requirePrivate(); err != nil {
		return nil, err
	}
	if arg.InstId == "" && arg.InstIdCode == 0 {
		return nil, errors.New("okx: ws amend order requires instId or instIdCode")
	}
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return nil, errors.New("okx: ws amend order requires ordId or clOrdId")
	}
	if arg.NewSz == "" && arg.NewPx == "" && arg.NewPxUsd == "" && arg.NewPxVol == "" {
		return nil, errors.New("okx: ws amend order requires newSz/newPx/newPxUsd/newPxVol")
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

func requiresPriceForOrderType(ordType string) bool {
	switch ordType {
	case "limit", "post_only", "fok", "ioc", "mmp", "mmp_and_post_only":
		return true
	default:
		return false
	}
}
