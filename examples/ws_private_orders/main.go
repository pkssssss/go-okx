package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx"
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
	instId := os.Getenv("OKX_INST_ID")

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

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()

	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:  okx.WSChannelOrders,
		InstType: instType,
		InstId:   instId,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("subscribed: channel=%s instType=%s instId=%s", okx.WSChannelOrders, instType, instId)

	unsubCtx, unsubCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer unsubCancel()

	if err := ws.UnsubscribeAndWait(unsubCtx, okx.WSArg{
		Channel:  okx.WSChannelOrders,
		InstType: instType,
		InstId:   instId,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("unsubscribed: channel=%s instType=%s instId=%s", okx.WSChannelOrders, instType, instId)
	ws.Close()
	<-ws.Done()
}
