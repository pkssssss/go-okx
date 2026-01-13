package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to purchase staking-defi; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	productId := os.Getenv("OKX_PRODUCT_ID")
	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_AMT")
	if productId == "" || ccy == "" || amt == "" {
		log.Fatal("missing env: OKX_PRODUCT_ID / OKX_CCY / OKX_AMT")
	}
	term := os.Getenv("OKX_TERM")
	tag := os.Getenv("OKX_TAG")
	demo := os.Getenv("OKX_DEMO") == "1"

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{APIKey: apiKey, SecretKey: secretKey, Passphrase: passphrase}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewFinanceStakingDefiPurchaseService().
		ProductId(productId).
		InvestData([]okx.FinanceStakingDefiInvestData{{Ccy: ccy, Amt: amt}})
	if term != "" {
		svc.Term(term)
	}
	if tag != "" {
		svc.Tag(tag)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ordId=%s tag=%s", ack.OrdId, ack.Tag)
}
