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

	bar := os.Getenv("OKX_CANDLE_BAR")
	if bar == "" {
		bar = "1m"
	}

	timeout := 10 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	channel := okx.WSCandleChannel(bar)
	if channel == "" {
		log.Fatal("invalid candle bar")
	}

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	outCh := make(chan okx.WSCandle, 1)
	ws.OnCandles(func(candle okx.WSCandle) {
		select {
		case outCh <- candle:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: channel, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case c := <-outCh:
		fmt.Printf("channel=%s instId=%s ts=%d o=%s h=%s l=%s c=%s vol=%s confirm=%s\n",
			c.Arg.Channel, c.Arg.InstId, c.Candle.TS, c.Candle.Open, c.Candle.High, c.Candle.Low, c.Candle.Close, c.Candle.Vol, c.Candle.Confirm)
		cancel()
	case <-time.After(timeout):
		log.Printf("timeout waiting message after %s", timeout)
		cancel()
	}

	<-ws.Done()
}
