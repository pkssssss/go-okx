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

	bal, err := c.NewAccountSubaccountBalancesService().SubAcct(subAcct).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_subaccount_balances: subAcct=%s uTime=%d totalEq=%s adjEq=%s details=%d", subAcct, bal.UTime, bal.TotalEq, bal.AdjEq, len(bal.Details))
	for i := 0; i < len(bal.Details) && i < 5; i++ {
		x := bal.Details[i]
		log.Printf("detail[%d]: ccy=%s eq=%s availBal=%s frozenBal=%s liab=%s", i, x.Ccy, x.Eq, x.AvailBal, x.FrozenBal, x.Liab)
	}
}
