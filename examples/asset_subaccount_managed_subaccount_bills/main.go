package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
	typ := os.Getenv("OKX_TYPE")
	subAcct := os.Getenv("OKX_SUB_ACCT")
	subUid := os.Getenv("OKX_SUB_UID")
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")

	var limit *int
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		limit = &n
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

	svc := c.NewAssetSubaccountManagedSubaccountBillsService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if typ != "" {
		svc.Type(typ)
	}
	if subAcct != "" {
		svc.SubAcct(subAcct)
	}
	if subUid != "" {
		svc.SubUid(subUid)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}
	if limit != nil {
		svc.Limit(*limit)
	}

	bills, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("managed_subaccount_bills=%d", len(bills))
	for i := 0; i < len(bills) && i < 5; i++ {
		x := bills[i]
		log.Printf("bill[%d]: billId=%s type=%s subAcct=%s subUid=%s ccy=%s amt=%s ts=%d", i, x.BillId, x.Type, x.SubAcct, x.SubUid, x.Ccy, x.Amt, x.TS)
	}
}
