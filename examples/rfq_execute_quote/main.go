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
		log.Fatal("refusing to execute quote; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	rfqId := os.Getenv("OKX_RFQ_ID")
	quoteId := os.Getenv("OKX_QUOTE_ID")
	if rfqId == "" || quoteId == "" {
		log.Fatal("missing env: OKX_RFQ_ID / OKX_QUOTE_ID")
	}

	var legs []okx.RFQExecuteQuoteLeg
	if v := os.Getenv("OKX_EXECUTE_QUOTE_LEGS"); v != "" {
		if err := json.Unmarshal([]byte(v), &legs); err != nil {
			log.Fatalf("invalid OKX_EXECUTE_QUOTE_LEGS: %v", err)
		}
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

	svc := c.NewRFQExecuteQuoteService().RfqId(rfqId).QuoteId(quoteId)
	if len(legs) > 0 {
		svc.Legs(legs)
	}

	trade, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("execute quote: blockTdId=%s isSuccessful=%v legs=%d", trade.BlockTdId, trade.IsSuccessful, len(trade.Legs))
}
