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

	msgCh := make(chan okx.WSSprdTrade, 1)
	errCh := make(chan error, 1)

	ws := c.NewWSBusinessPrivate(okx.WithWSTypedHandlerAsync(1024))
	ws.OnSprdTrades(func(trade okx.WSSprdTrade) {
		select {
		case msgCh <- trade:
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

	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelSprdTrades, SprdId: sprdId}); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed: channel=%s sprdId=%s", okx.WSChannelSprdTrades, sprdId)

	select {
	case tr := <-msgCh:
		log.Printf("sprd trade update: sprdId=%s tradeId=%s ordId=%s clOrdId=%s state=%s fillPx=%s fillSz=%s ts=%d legs=%d",
			tr.SprdId, tr.TradeId, tr.OrdId, tr.ClOrdId, tr.State, tr.FillPx, tr.FillSz, tr.TS, len(tr.Legs),
		)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no sprd trade update within %s (wait for filled/rejected updates)", timeout)
	}

	cancel()
	<-ws.Done()
}
