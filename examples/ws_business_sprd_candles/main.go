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
	bar := os.Getenv("OKX_BAR")
	if bar == "" {
		bar = "1D"
	}

	channel := okx.WSSprdCandleChannel(bar)
	if channel == "" {
		log.Fatal("empty candle channel")
	}

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	candleCh := make(chan okx.Candle, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseSprdCandles(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case candleCh <- dm.Data[0]:
		default:
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

	select {
	case c := <-candleCh:
		fmt.Printf("sprdId=%s channel=%s ts=%d o=%s h=%s l=%s c=%s vol=%s confirm=%s\n", sprdId, channel, c.TS, c.Open, c.High, c.Low, c.Close, c.Vol, c.Confirm)
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
