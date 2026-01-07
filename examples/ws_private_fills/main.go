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
		Channel: okx.WSChannelFills,
		InstId:  instId,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, func(message []byte) {
		ev, ok, err := okx.WSParseEvent(message)
		if err != nil || !ok {
			return
		}
		switch ev.Event {
		case "subscribe":
			if ev.Arg != nil && ev.Arg.Channel == okx.WSChannelFills {
				log.Printf("subscribed: channel=%s instId=%s connId=%s", ev.Arg.Channel, ev.Arg.InstId, ev.ConnID)
				cancel()
			}
		case "error":
			log.Printf("subscribe error: code=%s msg=%s", ev.Code, ev.Msg)
			cancel()
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting subscribe ack (fills channel requires VIP6+)")
		ws.Close()
		<-ws.Done()
	}
}
