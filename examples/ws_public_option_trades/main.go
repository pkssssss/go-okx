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
	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		instFamily = "BTC-USD"
	}
	instId := os.Getenv("OKX_INST_ID") // 可选：若传 instId，则以 instId 为主

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

	outCh := make(chan okx.WSOptionTrade, 1)
	ws.OnOptionTrades(func(trade okx.WSOptionTrade) {
		select {
		case outCh <- trade:
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

	arg := okx.WSArg{
		Channel:    okx.WSChannelOptionTrades,
		InstType:   "OPTION",
		InstId:     instId,
		InstFamily: instFamily,
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, arg); err != nil {
		log.Fatal(err)
	}

	select {
	case tr := <-outCh:
		fmt.Printf("instId=%s tradeId=%s side=%s px=%s sz=%s ts=%d\n", tr.InstId, tr.TradeId, tr.Side, tr.Px, tr.Sz, tr.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("no update within %s (option-trades is event-driven)", timeout)
		cancel()
	}

	<-ws.Done()
}
