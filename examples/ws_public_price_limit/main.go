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
		instId = "BTC-USDT-SWAP"
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

	outCh := make(chan okx.PriceLimit, 1)
	ws.OnPriceLimit(func(limit okx.PriceLimit) {
		select {
		case outCh <- limit:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelPriceLimit, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case pl := <-outCh:
		fmt.Printf("instId=%s buyLmt=%s sellLmt=%s enabled=%v ts=%d\n", pl.InstId, pl.BuyLmt, pl.SellLmt, pl.Enabled, pl.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("no update within %s (price-limit only pushes when changed)", timeout)
		cancel()
	}

	<-ws.Done()
}
