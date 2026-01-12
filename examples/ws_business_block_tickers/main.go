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

	outCh := make(chan okx.WSBlockTicker, 1)
	ws.OnBlockTickers(func(ticker okx.WSBlockTicker) {
		select {
		case outCh <- ticker:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelBlockTickers, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case tk := <-outCh:
		fmt.Printf("instType=%s instId=%s volCcy24h=%s vol24h=%s ts=%d\n", tk.InstType, tk.InstId, tk.VolCcy24h, tk.Vol24h, tk.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("timeout waiting message after %s", timeout)
		cancel()
	}

	<-ws.Done()
}
