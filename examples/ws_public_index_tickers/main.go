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

	c := okx.NewClient()
	ws := c.NewWSPublic()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	itCh := make(chan okx.IndexTicker, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseIndexTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case itCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelIndexTickers, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case it := <-itCh:
		fmt.Printf("instId=%s idxPx=%s open24h=%s high24h=%s low24h=%s ts=%d\n", it.InstId, it.IdxPx, it.Open24h, it.High24h, it.Low24h, it.TS)
		cancel()
	case <-time.After(70 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
