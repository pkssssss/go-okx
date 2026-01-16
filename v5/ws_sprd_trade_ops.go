package okx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	wsOpSprdOrder       = "sprd-order"
	wsOpSprdCancelOrder = "sprd-cancel-order"
	wsOpSprdAmendOrder  = "sprd-amend-order"
	wsOpSprdMassCancel  = "sprd-mass-cancel"
)

var errWSBusinessRequired = errors.New("okx: ws business client required")

func (w *WSClient) requireBusinessPrivate() error {
	if w == nil {
		return errors.New("okx: nil ws client")
	}
	if !w.needLogin {
		return errWSPrivateRequired
	}
	if w.kind != wsKindBusiness {
		return errWSBusinessRequired
	}
	return nil
}

// WSSprdPlaceOrderArg 表示 WS / 价差下单 的 args 项（精简版）。
// 数值字段保持为 string（无损）。
type WSSprdPlaceOrderArg struct {
	SprdId  string `json:"sprdId"`
	ClOrdId string `json:"clOrdId,omitempty"`
	Tag     string `json:"tag,omitempty"`

	Side    string `json:"side"`
	OrdType string `json:"ordType"`

	Px string `json:"px,omitempty"`
	Sz string `json:"sz"`
}

type WSSprdCancelOrderArg struct {
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
}

type WSSprdAmendOrderArg struct {
	OrdId   string `json:"ordId,omitempty"`
	ClOrdId string `json:"clOrdId,omitempty"`
	ReqId   string `json:"reqId,omitempty"`
	NewSz   string `json:"newSz,omitempty"`
	NewPx   string `json:"newPx,omitempty"`
}

type WSSprdMassCancelArg struct {
	SprdId string `json:"sprdId,omitempty"`
}

func validateWSSprdPlaceOrderArg(prefix string, arg WSSprdPlaceOrderArg) error {
	if arg.SprdId == "" {
		return fmt.Errorf("%s requires sprdId", prefix)
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
	if requiresPriceForOrderType(arg.OrdType) && arg.Px == "" {
		return fmt.Errorf("%s requires px for this ordType", prefix)
	}
	return nil
}

func validateWSSprdCancelOrderArg(prefix string, arg WSSprdCancelOrderArg) error {
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return fmt.Errorf("%s requires ordId or clOrdId", prefix)
	}
	return nil
}

func validateWSSprdAmendOrderArg(prefix string, arg WSSprdAmendOrderArg) error {
	if arg.OrdId == "" && arg.ClOrdId == "" {
		return fmt.Errorf("%s requires ordId or clOrdId", prefix)
	}
	if arg.NewSz == "" && arg.NewPx == "" {
		return fmt.Errorf("%s requires newSz or newPx", prefix)
	}
	return nil
}

// SprdPlaceOrder 通过 business WS 下单（op=sprd-order）。
// 注意：该方法会真实提交订单（建议在模拟盘验证，且调用方务必设置超时 ctx）。
func (w *WSClient) SprdPlaceOrder(ctx context.Context, arg WSSprdPlaceOrderArg) (*TradeOrderAck, error) {
	if err := w.requireBusinessPrivate(); err != nil {
		return nil, err
	}
	if err := validateWSSprdPlaceOrderArg("okx: ws sprd place order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpSprdOrder, []WSSprdPlaceOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws sprd place order empty response data")
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

// SprdCancelOrder 通过 business WS 撤单（op=sprd-cancel-order）。
func (w *WSClient) SprdCancelOrder(ctx context.Context, arg WSSprdCancelOrderArg) (*TradeOrderAck, error) {
	if err := w.requireBusinessPrivate(); err != nil {
		return nil, err
	}
	if err := validateWSSprdCancelOrderArg("okx: ws sprd cancel order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpSprdCancelOrder, []WSSprdCancelOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws sprd cancel order empty response data")
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

// SprdAmendOrder 通过 business WS 改单（op=sprd-amend-order）。
func (w *WSClient) SprdAmendOrder(ctx context.Context, arg WSSprdAmendOrderArg) (*TradeOrderAck, error) {
	if err := w.requireBusinessPrivate(); err != nil {
		return nil, err
	}
	if err := validateWSSprdAmendOrderArg("okx: ws sprd amend order", arg); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpSprdAmendOrder, []WSSprdAmendOrderArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalTradeOrderAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws sprd amend order empty response data")
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

func unmarshalSprdMassCancelAcks(reply *WSOpReply, raw []byte) ([]SprdMassCancelAck, error) {
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

	var acks []SprdMassCancelAck
	if err := json.Unmarshal(reply.Data, &acks); err != nil {
		return nil, err
	}
	return acks, nil
}

// SprdMassCancel 通过 business WS 全部撤单（op=sprd-mass-cancel）。
func (w *WSClient) SprdMassCancel(ctx context.Context, arg WSSprdMassCancelArg) (*SprdMassCancelAck, error) {
	if err := w.requireBusinessPrivate(); err != nil {
		return nil, err
	}

	reply, raw, err := w.doOpAndWaitRaw(ctx, wsOpSprdMassCancel, []WSSprdMassCancelArg{arg})
	if err != nil {
		return nil, err
	}
	acks, err := unmarshalSprdMassCancelAcks(reply, raw)
	if err != nil {
		return nil, err
	}
	if len(acks) == 0 {
		return nil, errors.New("okx: ws sprd mass cancel empty response data")
	}
	return &acks[0], nil
}
