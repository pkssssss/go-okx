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
	_ = ws.Subscribe(okx.WSArg{
		Channel:  okx.WSChannelOrders,
		InstType: instType,
		InstId:   instId,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, func(message []byte) {
		if ev, ok, err := okx.WSParseEvent(message); err == nil && ok {
			if ev.Event == "subscribe" && ev.Arg != nil && ev.Arg.Channel == okx.WSChannelOrders {
				log.Printf("subscribed: channel=%s instType=%s instId=%s connId=%s", ev.Arg.Channel, ev.Arg.InstType, ev.Arg.InstId, ev.ConnID)
				cancel()
			}
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting subscribe ack")
		ws.Close()
		<-ws.Done()
	}
}
