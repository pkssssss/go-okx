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

	ccy := os.Getenv("OKX_CCY") // optional; supports multi, comma separated
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

	svc := c.NewAccountSubaccountMaxWithdrawalService().SubAcct(subAcct)
	if ccy != "" {
		svc.Ccy(ccy)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_subaccount_max_withdrawal: subAcct=%s ccy=%s items=%d", subAcct, ccy, len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: ccy=%s maxWd=%s maxWdEx=%s spotOffsetMaxWd=%s spotOffsetMaxWdEx=%s", i, x.Ccy, x.MaxWd, x.MaxWdEx, x.SpotOffsetMaxWd, x.SpotOffsetMaxWdEx)
	}
}
