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
	instFamily := os.Getenv("OKX_INST_FAMILY")

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

	orderCh := make(chan okx.TradeOrder, 1)
	errCh := make(chan error, 1)

	ws := c.NewWSPrivate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseOrders(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case orderCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
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

	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:    okx.WSChannelOrders,
		InstType:   instType,
		InstId:     instId,
		InstFamily: instFamily,
	}); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed: channel=orders instType=%s instId=%s instFamily=%s", instType, instId, instFamily)

	select {
	case o := <-orderCh:
		log.Printf("order update: instId=%s ordId=%s side=%s ordType=%s px=%s sz=%s state=%s accFillSz=%s avgPx=%s uTime=%d",
			o.InstId, o.OrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State, o.AccFillSz, o.AvgPx, o.UTime,
		)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no order update within %s (place/cancel/amend an order to trigger push)", timeout)
	}

	cancel()
	<-ws.Done()
}
