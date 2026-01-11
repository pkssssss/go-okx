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
		log.Fatal("refusing to create quote; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	rfqId := os.Getenv("OKX_RFQ_ID")
	quoteSide := os.Getenv("OKX_QUOTE_SIDE")
	legsJSON := os.Getenv("OKX_QUOTE_LEGS") // JSON array
	if rfqId == "" || quoteSide == "" || legsJSON == "" {
		log.Fatal("missing env: OKX_RFQ_ID / OKX_QUOTE_SIDE / OKX_QUOTE_LEGS (JSON array)")
	}

	var legs []okx.QuoteLeg
	if err := json.Unmarshal([]byte(legsJSON), &legs); err != nil {
		log.Fatalf("invalid OKX_QUOTE_LEGS: %v", err)
	}

	anonymous, hasAnonymous := os.LookupEnv("OKX_QUOTE_ANONYMOUS")

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

	svc := c.NewRFQCreateQuoteService().
		RfqId(rfqId).
		ClQuoteId(os.Getenv("OKX_CL_QUOTE_ID")).
		Tag(os.Getenv("OKX_TAG")).
		QuoteSide(quoteSide).
		ExpiresIn(os.Getenv("OKX_QUOTE_EXPIRES_IN")).
		Legs(legs)
	if hasAnonymous {
		svc.Anonymous(anonymous == "true" || anonymous == "1")
	}

	quote, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("create quote: quoteId=%s rfqId=%s state=%s legs=%d", quote.QuoteId, quote.RfqId, quote.State, len(quote.Legs))
}
