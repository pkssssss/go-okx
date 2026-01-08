package main

import (
	"context"
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
	extraParams := os.Getenv("OKX_WS_EXTRA_PARAMS")

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

	ws.OnPositions(func(position okx.AccountPosition) {
		log.Printf("position push: instType=%s instId=%s posSide=%s pos=%s avgPx=%s uTime=%d", position.InstType, position.InstId, position.PosSide, position.Pos, position.AvgPx, position.UTime)
		cancel()
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:     okx.WSChannelPositions,
		InstType:    instType,
		InstFamily:  instFamily,
		InstId:      instId,
		ExtraParams: extraParams,
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting positions push")
		ws.Close()
		<-ws.Done()
	}
}
