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
		instId = "ONDO-USDC"
	}

	timeout := 10 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	outCh := make(chan okx.WSCallAuctionDetails, 1)
	ws.OnCallAuctionDetails(func(detail okx.WSCallAuctionDetails) {
		select {
		case outCh <- detail:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelCallAuctionDetails, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case d := <-outCh:
		fmt.Printf("instId=%s eqPx=%s matchedSz=%s unmatchedSz=%s state=%s auctionEndTime=%d ts=%d\n",
			d.InstId, d.EqPx, d.MatchedSz, d.UnmatchedSz, d.State, d.AuctionEndTime, d.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("no update within %s (call-auction-details is event-driven)", timeout)
		cancel()
	}

	<-ws.Done()
}
