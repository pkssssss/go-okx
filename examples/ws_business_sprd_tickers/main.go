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

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickerCh := make(chan okx.MarketSprdTicker, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseSprdTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case tickerCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelSprdTickers, SprdId: sprdId}); err != nil {
		log.Fatal(err)
	}

	select {
	case t := <-tickerCh:
		fmt.Printf("sprdId=%s last=%s lastSz=%s bid=%s/%s ask=%s/%s ts=%d\n", t.SprdId, t.Last, t.LastSz, t.BidPx, t.BidSz, t.AskPx, t.AskSz, t.TS)
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
