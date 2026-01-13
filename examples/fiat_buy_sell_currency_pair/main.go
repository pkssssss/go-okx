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

	pairs, err := c.NewFiatBuySellCurrencyPairService().
		FromCcy(fromCcy).
		ToCcy(toCcy).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fiat_buy_sell_currency_pair: n=%d from=%s to=%s", len(pairs), fromCcy, toCcy)
	for i, p := range pairs {
		log.Printf("pair[%d]: side=%s min=%s max=%s paymentMethods=%v", i, p.Side, p.SingleTradeMin, p.SingleTradeMax, p.PaymentMethods)
	}
}
