package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	channel := os.Getenv("OKX_BOOKS_CHANNEL")
	if channel == "" {
		channel = okx.WSChannelBooks
	}

	timeout := 10 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	store := okx.NewWSOrderBookStore(channel, instId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var applied int
	ws.OnOrderBook(func(dm okx.WSData[okx.WSOrderBook]) {
		if err := store.Apply(&dm); err != nil {
			log.Printf("order book error: %v", err)
			cancel()
			return
		}
		if !store.Ready() {
			return
		}

		applied++
		snap := store.Snapshot()
		if len(snap.Bids) == 0 || len(snap.Asks) == 0 {
			return
		}

		bestBid := snap.Bids[0]
		bestAsk := snap.Asks[0]
		fmt.Printf("channel=%s instId=%s seqId=%d ts=%d bid=%s/%s ask=%s/%s checksum=%d action=%s\n",
			snap.Channel, snap.InstId, snap.SeqId, snap.TS, bestBid.Px, bestBid.Sz, bestAsk.Px, bestAsk.Sz, snap.Checksum, dm.Action)

		// books 系列：通常第一条为 snapshot，第二条为 update；示例收到两条即退出。
		if applied >= 2 {
			cancel()
		}
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: channel, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(timeout):
		log.Printf("timeout after %s", timeout)
		ws.Close()
		<-ws.Done()
	}
}
