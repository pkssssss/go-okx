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
		instType = "FUTURES"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	instId := os.Getenv("OKX_INST_ID")
	if instFamily == "" && instId == "" {
		instFamily = "BTC-USD"
	}

	timeout := 30 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	msgCh := make(chan okx.EstimatedPrice, 1)
	ws.OnEstimatedPrice(func(price okx.EstimatedPrice) {
		select {
		case msgCh <- price:
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

	subArg := okx.WSArg{Channel: okx.WSChannelEstimatedPrice, InstType: instType}
	if instId != "" {
		subArg.InstId = instId
	} else {
		subArg.InstFamily = instFamily
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, subArg); err != nil {
		log.Fatal(err)
	}

	select {
	case price := <-msgCh:
		fmt.Printf("instType=%s instId=%s settlePx=%s settleType=%s ts=%d\n",
			price.InstType, price.InstId, price.SettlePx, price.SettleType, price.TS)
	case <-time.After(timeout):
		log.Printf("no estimated-price push within %s (only pushes ~1h before settlement/exercise)", timeout)
	}

	cancel()
	<-ws.Done()
}
