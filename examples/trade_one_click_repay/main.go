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
		log.Fatal("refusing to one-click repay; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	debtCcyListJSON := os.Getenv("OKX_DEBT_CCY_LIST")
	repayCcy := os.Getenv("OKX_REPAY_CCY")
	if debtCcyListJSON == "" || repayCcy == "" {
		log.Fatal("missing env: OKX_DEBT_CCY_LIST (JSON array) / OKX_REPAY_CCY")
	}

	var debtCcy []string
	if err := json.Unmarshal([]byte(debtCcyListJSON), &debtCcy); err != nil {
		log.Fatalf("invalid OKX_DEBT_CCY_LIST: %v", err)
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

	acks, err := c.NewOneClickRepayService().DebtCcy(debtCcy).RepayCcy(repayCcy).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_one_click_repay: acks=%d", len(acks))
	for i := 0; i < len(acks) && i < 5; i++ {
		x := acks[i]
		log.Printf("ack[%d]: debtCcy=%s repayCcy=%s status=%s uTime=%d", i, x.DebtCcy, x.RepayCcy, x.Status, x.UTime)
	}
}
