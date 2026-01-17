package okx

import (
	"testing"
	"time"
)

func TestWSArgKey_UniqueByUID(t *testing.T) {
	c := NewClient()
	w := c.NewWSPublic()

	a1 := WSArg{Channel: WSChannelDepositInfo, UID: "u1"}
	a2 := WSArg{Channel: WSChannelDepositInfo, UID: "u2"}

	waiter := w.registerWaiter("1", "subscribe", []WSArg{a1, a2})

	w.notifyWaiter(WSEvent{ID: "1", Event: "subscribe", Arg: &a1})
	select {
	case <-waiter.done:
		t.Fatalf("waiter done too early")
	default:
	}

	w.notifyWaiter(WSEvent{ID: "1", Event: "subscribe", Arg: &a2})
	select {
	case err := <-waiter.done:
		if err != nil {
			t.Fatalf("waiter err = %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timeout waiting waiter done")
	}
}

func TestWSArgKey_UniqueByExtraParams(t *testing.T) {
	c := NewClient()
	w := c.NewWSPublic()

	a1 := WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT", ExtraParams: "x=1"}
	a2 := WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT", ExtraParams: "x=2"}

	waiter := w.registerWaiter("1", "subscribe", []WSArg{a1, a2})

	w.notifyWaiter(WSEvent{ID: "1", Event: "subscribe", Arg: &a1})
	select {
	case <-waiter.done:
		t.Fatalf("waiter done too early")
	default:
	}

	w.notifyWaiter(WSEvent{ID: "1", Event: "subscribe", Arg: &a2})
	select {
	case err := <-waiter.done:
		if err != nil {
			t.Fatalf("waiter err = %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timeout waiting waiter done")
	}
}
