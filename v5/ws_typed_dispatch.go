package okx

import (
	"context"
	"fmt"
)

type wsTypedKind int

const (
	wsTypedKindOrders wsTypedKind = iota + 1
	wsTypedKindFills
	wsTypedKindAccount
	wsTypedKindPositions
	wsTypedKindBalanceAndPosition
	wsTypedKindDepositInfo
	wsTypedKindWithdrawalInfo
	wsTypedKindOpReply
)

func (k wsTypedKind) String() string {
	switch k {
	case wsTypedKindOrders:
		return "orders"
	case wsTypedKindFills:
		return "fills"
	case wsTypedKindAccount:
		return "account"
	case wsTypedKindPositions:
		return "positions"
	case wsTypedKindBalanceAndPosition:
		return "balance_and_position"
	case wsTypedKindDepositInfo:
		return "deposit_info"
	case wsTypedKindWithdrawalInfo:
		return "withdrawal_info"
	case wsTypedKindOpReply:
		return "op_reply"
	default:
		return "unknown"
	}
}

type wsTypedTask struct {
	kind wsTypedKind

	orders         []TradeOrder
	fills          []WSFill
	balances       []AccountBalance
	positions      []AccountPosition
	balPos         []WSBalanceAndPosition
	depositInfo    []WSDepositInfo
	withdrawalInfo []WSWithdrawalInfo

	op    WSOpReply
	opRaw []byte
}

func (w *WSClient) typedDispatchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-w.typedQueue:
			w.handleTyped(task)
		}
	}
}

func (w *WSClient) dispatchTyped(task wsTypedTask) {
	if w == nil {
		return
	}

	if !w.typedAsync || w.typedQueue == nil {
		w.handleTyped(task)
		return
	}

	select {
	case w.typedQueue <- task:
	default:
		w.onError(fmt.Errorf("okx: ws typed handler queue full; dropped kind=%s", task.kind.String()))
	}
}

func (w *WSClient) handleTyped(task wsTypedTask) {
	if w == nil {
		return
	}

	switch task.kind {
	case wsTypedKindOrders:
		w.typedMu.RLock()
		h := w.ordersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.orders) == 0 {
			return
		}
		for _, order := range task.orders {
			o := order
			w.safeTypedCall(task.kind, func() { h(o) })
		}
	case wsTypedKindFills:
		w.typedMu.RLock()
		h := w.fillsHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.fills) == 0 {
			return
		}
		for _, fill := range task.fills {
			f := fill
			w.safeTypedCall(task.kind, func() { h(f) })
		}
	case wsTypedKindAccount:
		w.typedMu.RLock()
		h := w.accountHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.balances) == 0 {
			return
		}
		for _, balance := range task.balances {
			b := balance
			w.safeTypedCall(task.kind, func() { h(b) })
		}
	case wsTypedKindPositions:
		w.typedMu.RLock()
		h := w.positionsHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.positions) == 0 {
			return
		}
		for _, position := range task.positions {
			p := position
			w.safeTypedCall(task.kind, func() { h(p) })
		}
	case wsTypedKindBalanceAndPosition:
		w.typedMu.RLock()
		h := w.balanceAndPositionHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.balPos) == 0 {
			return
		}
		for _, data := range task.balPos {
			d := data
			w.safeTypedCall(task.kind, func() { h(d) })
		}
	case wsTypedKindDepositInfo:
		w.typedMu.RLock()
		h := w.depositInfoHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.depositInfo) == 0 {
			return
		}
		for _, info := range task.depositInfo {
			i := info
			w.safeTypedCall(task.kind, func() { h(i) })
		}
	case wsTypedKindWithdrawalInfo:
		w.typedMu.RLock()
		h := w.withdrawalInfoHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.withdrawalInfo) == 0 {
			return
		}
		for _, info := range task.withdrawalInfo {
			i := info
			w.safeTypedCall(task.kind, func() { h(i) })
		}
	case wsTypedKindOpReply:
		w.typedMu.RLock()
		h := w.opReplyHandler
		w.typedMu.RUnlock()
		if h == nil {
			return
		}
		w.safeTypedCall(task.kind, func() { h(task.op, task.opRaw) })
	default:
		return
	}
}

func (w *WSClient) safeTypedCall(kind wsTypedKind, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			w.onError(fmt.Errorf("okx: ws typed handler panic kind=%s: %v", kind.String(), r))
		}
	}()
	fn()
}
