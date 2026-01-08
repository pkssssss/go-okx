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

	c := okx.NewClient()
	ws := c.NewWSPublic()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	osCh := make(chan okx.OptSummary, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseOptSummary(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case osCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelOptSummary, InstFamily: instFamily}); err != nil {
		log.Fatal(err)
	}

	select {
	case s := <-osCh:
		fmt.Printf("instId=%s uly=%s delta=%s gamma=%s theta=%s vega=%s ts=%d\n", s.InstId, s.Uly, s.Delta, s.Gamma, s.Theta, s.Vega, s.TS)
		cancel()
	case <-time.After(30 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
