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

	tkCh := make(chan okx.MarketTicker, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case tkCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelTickers, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case tk := <-tkCh:
		fmt.Printf("instId=%s last=%s bid=%s/%s ask=%s/%s ts=%d\n", tk.InstId, tk.Last, tk.BidPx, tk.BidSz, tk.AskPx, tk.AskSz, tk.TS)
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
