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
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	outCh := make(chan okx.BlockTrade, 1)
	ws.OnPublicBlockTrades(func(trade okx.BlockTrade) {
		select {
		case outCh <- trade:
		default:
		}
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelPublicBlockTrades, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case tr := <-outCh:
		fmt.Printf("instId=%s tradeId=%s px=%s sz=%s side=%s groupId=%s ts=%d\n", tr.InstId, tr.TradeId, tr.Px, tr.Sz, tr.Side, tr.GroupId, tr.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("timeout waiting message after %s", timeout)
		cancel()
	}

	<-ws.Done()
}
