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
	ccy := os.Getenv("OKX_CCY")

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

	msgCh := make(chan okx.WSWithdrawalInfo, 1)
	errCh := make(chan error, 1)

	ws := c.NewWSBusinessPrivate(okx.WithWSTypedHandlerAsync(1024))
	ws.OnWithdrawalInfo(func(info okx.WSWithdrawalInfo) {
		select {
		case msgCh <- info:
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

	subArg := okx.WSArg{Channel: okx.WSChannelWithdrawalInfo}
	if ccy != "" {
		subArg.Ccy = ccy
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()

	if err := ws.SubscribeAndWait(subCtx, subArg); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed: channel=%s ccy=%s", okx.WSChannelWithdrawalInfo, ccy)

	select {
	case m := <-msgCh:
		log.Printf("withdrawal update: ccy=%s chain=%s wdId=%s amt=%s state=%s ts=%d pTime=%d subAcct=%s uid=%s fee=%s feeCcy=%s txId=%s",
			m.Ccy, m.Chain, m.WdId, m.Amt, m.State, m.TS, m.PTime, m.SubAcct, m.UID, m.Fee, m.FeeCcy, m.TxId,
		)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no withdrawal update within %s (withdraw/cancel or wait for state change)", timeout)
	}

	cancel()
	<-ws.Done()
}
