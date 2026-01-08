package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tradeCh := make(chan okx.MarketTrade, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseTradesAll(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case tradeCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelTradesAll, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case t := <-tradeCh:
		fmt.Printf("instId=%s tradeId=%s px=%s sz=%s side=%s ts=%d source=%s\n", t.InstId, t.TradeId, t.Px, t.Sz, t.Side, t.TS, t.Source)
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
