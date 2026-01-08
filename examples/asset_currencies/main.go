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

	demo := os.Getenv("OKX_DEMO") == "1"
	ccy := os.Getenv("OKX_CCY")

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

	svc := c.NewAssetCurrenciesService()
	if ccy != "" {
		svc.Ccy(ccy)
	}

	ccys, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("currencies=%d", len(ccys))
	for i := 0; i < len(ccys) && i < 5; i++ {
		x := ccys[i]
		log.Printf("ccy[%d]: ccy=%s chain=%s canDep=%v canWd=%v fee=%s wdTickSz=%s needTag=%v mainNet=%v", i, x.Ccy, x.Chain, x.CanDep, x.CanWd, x.Fee, x.WdTickSz, x.NeedTag, x.MainNet)
	}
}
