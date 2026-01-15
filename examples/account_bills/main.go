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

	svc := c.NewAccountBillsService()
	if v := os.Getenv("OKX_INST_TYPE"); v != "" {
		svc.InstType(v)
	}
	if v := os.Getenv("OKX_INST_ID"); v != "" {
		svc.InstId(v)
	}
	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}
	if v := os.Getenv("OKX_MGN_MODE"); v != "" {
		svc.MgnMode(v)
	}
	if v := os.Getenv("OKX_CT_TYPE"); v != "" {
		svc.CtType(v)
	}
	if v := os.Getenv("OKX_TYPE"); v != "" {
		svc.Type(v)
	}
	if v := os.Getenv("OKX_SUB_TYPE"); v != "" {
		svc.SubType(v)
	}
	if v := os.Getenv("OKX_AFTER"); v != "" {
		svc.After(v)
	}
	if v := os.Getenv("OKX_BEFORE"); v != "" {
		svc.Before(v)
	}
	if v := os.Getenv("OKX_BEGIN"); v != "" {
		svc.Begin(v)
	}
	if v := os.Getenv("OKX_END"); v != "" {
		svc.End(v)
	}
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid env OKX_LIMIT=%q: %v", v, err)
		}
		svc.Limit(n)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_bills: items=%d", len(items))
	if len(items) > 0 {
		x := items[0]
		log.Printf("first: billId=%s instType=%s instId=%s type=%s subType=%s ccy=%s balChg=%s ts=%d", x.BillId, x.InstType, x.InstId, x.Type, x.SubType, x.Ccy, x.BalChg, x.TS)
	}
}
