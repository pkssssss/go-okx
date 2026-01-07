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

	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseBalanceAndPosition(message)
		if err != nil || !ok {
			return
		}
		if len(dm.Data) == 0 {
			log.Printf("balance_and_position push: empty data")
			cancel()
			return
		}

		d := dm.Data[0]
		log.Printf("balance_and_position push: pTime=%d eventType=%s bal=%d pos=%d", d.PTime, d.EventType, len(d.BalData), len(d.PosData))
		cancel()
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelBalanceAndPosition}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting balance_and_position push")
		ws.Close()
		<-ws.Done()
	}
}
