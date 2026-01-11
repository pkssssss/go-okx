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
	ccy := os.Getenv("OKX_CCY")

	timeout := 15 * time.Second
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

	gotCh := make(chan okx.AccountGreeks, 1)
	ws.OnAccountGreeks(func(greeks okx.AccountGreeks) {
		select {
		case gotCh <- greeks:
		default:
		}
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelAccountGreeks, Ccy: ccy}); err != nil {
		log.Fatal(err)
	}

	select {
	case greeks := <-gotCh:
		fmt.Printf("ccy=%s deltaBS=%s deltaPA=%s gammaBS=%s gammaPA=%s thetaBS=%s thetaPA=%s vegaBS=%s vegaPA=%s ts=%d\n",
			greeks.Ccy, greeks.DeltaBS, greeks.DeltaPA, greeks.GammaBS, greeks.GammaPA, greeks.ThetaBS, greeks.ThetaPA, greeks.VegaBS, greeks.VegaPA, greeks.TS)
		cancel()
	case <-time.After(timeout):
		log.Printf("timeout after %s", timeout)
		ws.Close()
		<-ws.Done()
		return
	}

	<-ws.Done()
}
