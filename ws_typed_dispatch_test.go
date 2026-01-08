package okx

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestWSClient_DispatchTyped_DropsWhenQueueFull(t *testing.T) {
	errCh := make(chan error, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask), // unbuffered + no receiver => full
		errHandler: func(err error) {
			select {
			case errCh <- err:
			default:
			}
		},
	}

	w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrders, orders: []TradeOrder{{OrdId: "o1"}}})

	select {
	case err := <-errCh:
		if err == nil || !strings.Contains(err.Error(), "queue full") || !strings.Contains(err.Error(), "kind=orders") {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting error")
	}
}

func TestWSClient_TypedDispatchLoop_PanicRecovered(t *testing.T) {
	errCh := make(chan error, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
		errHandler: func(err error) {
			select {
			case errCh <- err:
			default:
			}
		},
	}

	w.OnOrders(func(order TradeOrder) {
		_ = order
		panic("boom")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrders, orders: []TradeOrder{{OrdId: "o1"}}})

	select {
	case err := <-errCh:
		if err == nil || !strings.Contains(err.Error(), "panic") || !strings.Contains(err.Error(), "kind=orders") {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting error")
	}
}

func TestWSClient_onDataMessage_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan TradeOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnOrders(func(order TradeOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"orders"},"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","ordId":"o1","clOrdId":"c1"}]}`))

	select {
	case o := <-gotCh:
		if o.OrdId != "o1" || o.ClOrdId != "c1" {
			t.Fatalf("order = %#v", o)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting order")
	}
}
