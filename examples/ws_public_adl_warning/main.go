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
		instType = "FUTURES"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		instFamily = "BTC-USDT"
	}

	timeout := 90 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	msgCh := make(chan okx.WSADLWarning, 1)
	ws.OnADLWarning(func(warning okx.WSADLWarning) {
		select {
		case msgCh <- warning:
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
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelADLWarning, InstType: instType, InstFamily: instFamily}); err != nil {
		log.Fatal(err)
	}

	select {
	case w := <-msgCh:
		fmt.Printf("instType=%s instFamily=%s ccy=%s state=%s bal=%s maxBal=%s maxBalTs=%d adlType=%s ts=%d\n",
			w.InstType, w.InstFamily, w.Ccy, w.State, w.Bal, w.MaxBal, w.MaxBalTS, w.ADLType, w.TS)
	case <-time.After(timeout):
		log.Printf("no adl-warning push within %s (normal state pushes once per minute)", timeout)
	}

	cancel()
	<-ws.Done()
}
