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
	sprdId := os.Getenv("OKX_SPRD_ID")
	if sprdId == "" {
		sprdId = "BTC-USDT_BTC-USDT-SWAP"
	}

	channel := os.Getenv("OKX_BOOKS_CHANNEL")
	if channel == "" {
		channel = okx.WSChannelSprdBooks5
	}

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	store := okx.NewWSSprdOrderBookStore(channel, sprdId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var applied int
	if err := ws.Start(ctx, func(message []byte) {
		ok, err := store.ApplyMessage(message)
		if err != nil {
			log.Printf("order book error: %v", err)
			cancel()
			return
		}
		if !ok {
			return
		}

		applied++
		snap := store.Snapshot()
		if len(snap.Bids) == 0 || len(snap.Asks) == 0 {
			return
		}

		bestBid := snap.Bids[0]
		bestAsk := snap.Asks[0]
		fmt.Printf("channel=%s sprdId=%s seqId=%d ts=%d bid=%s/%s ask=%s/%s checksum=%d\n",
			snap.Channel, snap.SprdId, snap.SeqId, snap.TS, bestBid.Px, bestBid.Sz, bestAsk.Px, bestAsk.Sz, snap.Checksum)

		// sprd-books*：通常第一条为 snapshot，后续为 update/定量推送；示例收到两条即退出。
		if applied >= 2 {
			cancel()
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: channel, SprdId: sprdId}); err != nil {
		log.Fatal(err)
	}

	<-ws.Done()
}
