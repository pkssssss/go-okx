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
	typ := os.Getenv("OKX_ASSET_BILLS_HISTORY_TYPE")
	clientId := os.Getenv("OKX_ASSET_BILLS_HISTORY_CLIENT_ID")
	after := os.Getenv("OKX_ASSET_BILLS_HISTORY_AFTER")
	before := os.Getenv("OKX_ASSET_BILLS_HISTORY_BEFORE")
	pagingType := os.Getenv("OKX_ASSET_BILLS_HISTORY_PAGING_TYPE")

	var limit *int
	if v := os.Getenv("OKX_ASSET_BILLS_HISTORY_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_ASSET_BILLS_HISTORY_LIMIT: %v", err)
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

	svc := c.NewAssetBillsHistoryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if typ != "" {
		svc.Type(typ)
	}
	if clientId != "" {
		svc.ClientId(clientId)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}
	if pagingType != "" {
		svc.PagingType(pagingType)
	}
	if limit != nil {
		svc.Limit(*limit)
	}

	bills, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("asset bills history=%d", len(bills))
	for i := 0; i < len(bills) && i < 5; i++ {
		b := bills[i]
		log.Printf("bill[%d]: billId=%s ccy=%s type=%s balChg=%s bal=%s ts=%d notes=%s", i, b.BillId, b.Ccy, b.Type, b.BalChg, b.Bal, b.TS, b.Notes)
	}
}
