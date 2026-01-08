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

	fromCcy := os.Getenv("OKX_FROM_CCY")
	toCcy := os.Getenv("OKX_TO_CCY")
	if fromCcy == "" || toCcy == "" {
		log.Fatal("missing env: OKX_FROM_CCY / OKX_TO_CCY")
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

	pairs, err := c.NewAssetConvertCurrencyPairService().FromCcy(fromCcy).ToCcy(toCcy).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("pairs=%d", len(pairs))
	for i := 0; i < len(pairs) && i < 5; i++ {
		x := pairs[i]
		log.Printf("pair[%d]: instId=%s base=%s min=%s max=%s quote=%s min=%s max=%s", i, x.InstId, x.BaseCcy, x.BaseCcyMin, x.BaseCcyMax, x.QuoteCcy, x.QuoteCcyMin, x.QuoteCcyMax)
	}
}
