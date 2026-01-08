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

	svc := c.NewAssetNonTradableAssetsService()
	if ccy != "" {
		svc.Ccy(ccy)
	}

	assets, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("non_tradable_assets=%d", len(assets))
	for i := 0; i < len(assets) && i < 5; i++ {
		x := assets[i]
		log.Printf("asset[%d]: ccy=%s chain=%s bal=%s canWd=%v fee=%s %s wdTickSz=%s needTag=%v", i, x.Ccy, x.Chain, x.Bal, x.CanWd, x.Fee, x.FeeCcy, x.WdTickSz, x.NeedTag)
	}
}
