package okx

import (
	"errors"
	"testing"
	"time"
)

func TestWSClient_notifyWaiter_NonBlockingWhenDoneFull(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPublic()

	arg := WSArg{Channel: WSChannelTickers, InstId: "BTC-USDT"}
	waiter := ws.registerWaiter("1", "subscribe", []WSArg{arg})

	// 模拟 closeConn/failWaiters 已先行写入一次（buffer=1 满）。
	waiter.done <- errors.New("first")

	done := make(chan struct{})
	go func() {
		ws.notifyWaiter(WSEvent{ID: "1", Event: "subscribe", Arg: &arg})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("notifyWaiter blocked when waiter.done is full")
	}
}

func TestWSClient_notifyOpWaiter_NonBlockingWhenDoneFull(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPrivate()

	waiter := ws.registerOpWaiter("1", "order")

	// 模拟 closeConn/failOpWaiters 已先行写入一次（buffer=1 满）。
	waiter.done <- wsOpRespResult{err: errors.New("first")}

	reply := WSOpReply{ID: "1", Op: "order", Code: "0", Msg: ""}

	done := make(chan struct{})
	go func() {
		ws.notifyOpWaiter(reply, []byte(`{}`))
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("notifyOpWaiter blocked when waiter.done is full")
	}
}

func TestWSClient_notifyOpWaiterError_NonBlockingWhenDoneFull(t *testing.T) {
	c := NewClient()
	ws := c.NewWSPrivate()

	waiter := ws.registerOpWaiter("1", "order")

	// 模拟 closeConn/failOpWaiters 已先行写入一次（buffer=1 满）。
	waiter.done <- wsOpRespResult{err: errors.New("first")}

	done := make(chan struct{})
	go func() {
		ws.notifyOpWaiterError(WSEvent{ID: "1", Event: "error", Code: "60012", Msg: "Invalid request"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("notifyOpWaiterError blocked when waiter.done is full")
	}
}
