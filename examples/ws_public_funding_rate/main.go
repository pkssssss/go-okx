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

	frCh := make(chan okx.FundingRate, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseFundingRate(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case frCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelFundingRate, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case fr := <-frCh:
		fmt.Printf("instId=%s fundingRate=%s fundingTime=%d nextFundingTime=%d ts=%d\n", fr.InstId, fr.FundingRate, fr.FundingTime, fr.NextFundingTime, fr.TS)
		cancel()
	case <-time.After(120 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
