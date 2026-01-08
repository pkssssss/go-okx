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

	c := okx.NewClient()
	ws := c.NewWSPublic()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	oiCh := make(chan okx.OpenInterest, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseOpenInterest(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case oiCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelOpenInterest, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case oi := <-oiCh:
		fmt.Printf("instId=%s instType=%s oi=%s oiCcy=%s oiUsd=%s ts=%d\n", oi.InstId, oi.InstType, oi.OI, oi.OICcy, oi.OIUsd, oi.TS)
		cancel()
	case <-time.After(15 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
