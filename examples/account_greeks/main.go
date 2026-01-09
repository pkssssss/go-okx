package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	ccy := os.Getenv("OKX_CCY") // optional
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

	svc := c.NewAccountGreeksService()
	if ccy != "" {
		svc.Ccy(ccy)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_greeks: items=%d", len(data))
	for _, it := range data {
		log.Printf("ccy=%s ts=%d deltaBS=%s deltaPA=%s gammaBS=%s gammaPA=%s thetaBS=%s thetaPA=%s vegaBS=%s vegaPA=%s",
			it.Ccy, it.TS, it.DeltaBS, it.DeltaPA, it.GammaBS, it.GammaPA, it.ThetaBS, it.ThetaPA, it.VegaBS, it.VegaPA,
		)
	}
}
