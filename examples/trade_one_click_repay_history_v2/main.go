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

	svc := c.NewOneClickRepayHistoryV2Service()
	if v := os.Getenv("OKX_AFTER"); v != "" {
		svc.After(v)
	}
	if v := os.Getenv("OKX_BEFORE"); v != "" {
		svc.Before(v)
	}
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_one_click_repay_history_v2: items=%d", len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: debtCcy=%s fillDebtSz=%s status=%s ts=%d ordIdInfo=%d", i, x.DebtCcy, x.FillDebtSz, x.Status, x.TS, len(x.OrdIdInfo))
	}
}
