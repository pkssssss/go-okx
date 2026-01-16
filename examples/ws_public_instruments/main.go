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
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SPOT"
	}

	timeout := 30 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	msgCh := make(chan okx.Instrument, 1)
	ws.OnInstruments(func(instrument okx.Instrument) {
		select {
		case msgCh <- instrument:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelInstruments, InstType: instType}); err != nil {
		log.Fatal(err)
	}

	select {
	case inst := <-msgCh:
		fmt.Printf("instType=%s instId=%s base=%s quote=%s tickSz=%s lotSz=%s minSz=%s state=%s\n",
			inst.InstType, inst.InstId, inst.BaseCcy, inst.QuoteCcy, inst.TickSz, inst.LotSz, inst.MinSz, inst.State)
	case <-time.After(timeout):
		log.Printf("no instruments update within %s (channel is incremental; try later)", timeout)
	}

	cancel()
	<-ws.Done()
}
