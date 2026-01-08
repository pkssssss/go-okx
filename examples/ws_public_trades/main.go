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

	tradesCh := make(chan []okx.MarketTrade, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseTrades(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case tradesCh <- dm.Data:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelTrades, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case trades := <-tradesCh:
		t := trades[0]
		fmt.Printf("instId=%s tradeId=%s px=%s sz=%s side=%s ts=%d count=%s source=%s seqId=%d n=%d\n", t.InstId, t.TradeId, t.Px, t.Sz, t.Side, t.TS, t.Count, t.Source, t.SeqId, len(trades))
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
