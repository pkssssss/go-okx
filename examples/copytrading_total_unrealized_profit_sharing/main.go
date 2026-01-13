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

	instType := os.Getenv("OKX_INST_TYPE")

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

	svc := c.NewCopyTradingTotalUnrealizedProfitSharingService()
	if instType != "" {
		svc.InstType(instType)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_total_unrealized_profit_sharing: items=%d", len(items))
	for i := 0; i < len(items); i++ {
		it := items[i]
		log.Printf("item[%d]: profitSharingTs=%d totalUnrealizedProfitSharingAmt=%s", i, it.ProfitSharingTs, it.TotalUnrealizedProfitSharingAmt)
	}
}
