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

	sprdId := os.Getenv("OKX_SPRD_ID")
	if sprdId == "" {
		sprdId = "BTC-USDT_BTC-USDT-SWAP"
	}

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

	msgCh := make(chan okx.SprdOrder, 1)
	errCh := make(chan error, 1)

	ws := c.NewWSBusinessPrivate(okx.WithWSTypedHandlerAsync(1024))
	ws.OnSprdOrders(func(order okx.SprdOrder) {
		select {
		case msgCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}); err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()

	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelSprdOrders, SprdId: sprdId}); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed: channel=%s sprdId=%s", okx.WSChannelSprdOrders, sprdId)

	select {
	case o := <-msgCh:
		log.Printf("sprd order update: sprdId=%s ordId=%s clOrdId=%s side=%s ordType=%s px=%s sz=%s state=%s accFillSz=%s uTime=%d",
			o.SprdId, o.OrdId, o.ClOrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State, o.AccFillSz, o.UTime,
		)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no sprd order update within %s (place/cancel/amend to trigger)", timeout)
	}

	cancel()
	<-ws.Done()
}
