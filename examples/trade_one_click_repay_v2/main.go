package main

import (
	"context"
	"encoding/json"
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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to one-click repay v2; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	debtCcy := os.Getenv("OKX_DEBT_CCY")
	repayCcyListJSON := os.Getenv("OKX_REPAY_CCY_LIST")
	if debtCcy == "" || repayCcyListJSON == "" {
		log.Fatal("missing env: OKX_DEBT_CCY / OKX_REPAY_CCY_LIST (JSON array)")
	}

	var repayCcyList []string
	if err := json.Unmarshal([]byte(repayCcyListJSON), &repayCcyList); err != nil {
		log.Fatalf("invalid OKX_REPAY_CCY_LIST: %v", err)
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

	ack, err := c.NewOneClickRepayV2Service().DebtCcy(debtCcy).RepayCcyList(repayCcyList).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_one_click_repay_v2: debtCcy=%s repayCcyList=%v ts=%d", ack.DebtCcy, ack.RepayCcyList, ack.TS)
}
