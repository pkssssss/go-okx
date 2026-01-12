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

	outCh := make(chan okx.WSPublicStrucBlockTrade, 1)
	ws.OnPublicStrucBlockTrades(func(trade okx.WSPublicStrucBlockTrade) {
		select {
		case outCh <- trade:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelPublicStrucBlockTrades}); err != nil {
		log.Fatal(err)
	}

	select {
	case tr := <-outCh:
		leg := ""
		if len(tr.Legs) > 0 {
			leg = fmt.Sprintf(" firstLeg(instId=%s px=%s sz=%s side=%s tradeId=%s)", tr.Legs[0].InstId, tr.Legs[0].Px, tr.Legs[0].Sz, tr.Legs[0].Side, tr.Legs[0].TradeId)
		}
		fmt.Printf("blockTdId=%s groupId=%s cTime=%d%s\n", tr.BlockTdId, tr.GroupId, tr.CTime, leg)
		cancel()
	case <-time.After(timeout):
		log.Printf("timeout waiting message after %s", timeout)
		cancel()
	}

	<-ws.Done()
}
