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
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "ANY"
	}
	instFamily := os.Getenv("OKX_INST_FAMILY")
	instId := os.Getenv("OKX_INST_ID")

	timeout := 60 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{
			APIKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	ws := c.NewWSPrivate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws.OnLiquidationWarning(func(warning okx.WSLiquidationWarning) {
		fmt.Printf("instType=%s instId=%s mgnMode=%s posId=%s posSide=%s pos=%s lever=%s markPx=%s mgnRatio=%s ccy=%s pTime=%d\n",
			warning.InstType, warning.InstId, warning.MgnMode, warning.PosId, warning.PosSide, warning.Pos, warning.Lever, warning.MarkPx, warning.MgnRatio, warning.Ccy, warning.PTime)
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:    okx.WSChannelLiquidationWarning,
		InstType:   instType,
		InstFamily: instFamily,
		InstId:     instId,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("subscribed: channel=%s instType=%s instFamily=%s instId=%s", okx.WSChannelLiquidationWarning, instType, instFamily, instId)

	select {
	case <-ws.Done():
	case <-time.After(timeout):
		log.Printf("timeout after %s", timeout)
		ws.Close()
		<-ws.Done()
	}
}
