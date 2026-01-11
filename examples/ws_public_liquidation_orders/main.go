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
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
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

	outCh := make(chan okx.LiquidationOrder, 1)
	ws.OnLiquidationOrders(func(order okx.LiquidationOrder) {
		select {
		case outCh <- order:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelLiquidationOrders, InstType: instType}); err != nil {
		log.Fatal(err)
	}

	select {
	case lo := <-outCh:
		var d okx.LiquidationOrderDetail
		if len(lo.Details) > 0 {
			d = lo.Details[0]
		}
		fmt.Printf("instType=%s instId=%s uly=%s posSide=%s side=%s bkPx=%s sz=%s ts=%d\n",
			lo.InstType, lo.InstId, lo.Uly, d.PosSide, d.Side, d.BkPx, d.Sz, d.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("no update within %s (liquidation-orders is event-driven)", timeout)
		cancel()
	}

	<-ws.Done()
}
