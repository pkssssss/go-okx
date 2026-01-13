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

	instType := os.Getenv("OKX_INST_TYPE")
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")

	var limit *int
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		limit = &n
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

	svc := c.NewCopyTradingProfitSharingDetailsService()
	if instType != "" {
		svc.InstType(instType)
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

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_profit_sharing_details: items=%d", len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		it := items[i]
		log.Printf("item[%d]: instType=%s ccy=%s profitSharingAmt=%s profitSharingId=%s ts=%d nickName=%s", i, it.InstType, it.Ccy, it.ProfitSharingAmt, it.ProfitSharingId, it.TS, it.NickName)
	}
}
