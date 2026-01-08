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

	subAcct := os.Getenv("OKX_SUB_ACCT")
	if subAcct == "" {
		log.Fatal("missing env: OKX_SUB_ACCT")
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

	svc := c.NewAssetSubaccountBalancesService().SubAcct(subAcct)
	if ccy != "" {
		svc.Ccy(ccy)
	}

	bals, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("subaccount_balances=%d subAcct=%s", len(bals), subAcct)
	for i := 0; i < len(bals) && i < 5; i++ {
		x := bals[i]
		log.Printf("bal[%d]: ccy=%s bal=%s availBal=%s frozenBal=%s", i, x.Ccy, x.Bal, x.AvailBal, x.FrozenBal)
	}
}
