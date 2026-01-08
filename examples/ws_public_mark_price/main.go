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

	mpCh := make(chan okx.MarkPrice, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseMarkPrice(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case mpCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelMarkPrice, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case mp := <-mpCh:
		fmt.Printf("instId=%s instType=%s markPx=%s ts=%d\n", mp.InstId, mp.InstType, mp.MarkPx, mp.TS)
		cancel()
	case <-time.After(15 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
