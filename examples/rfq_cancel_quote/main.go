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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel quote; set OKX_CONFIRM=YES to continue")
	}

	quoteId := os.Getenv("OKX_QUOTE_ID")
	clQuoteId := os.Getenv("OKX_CL_QUOTE_ID")
	if quoteId == "" && clQuoteId == "" {
		log.Fatal("missing env: OKX_QUOTE_ID or OKX_CL_QUOTE_ID")
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

	svc := c.NewRFQCancelQuoteService().
		RfqId(os.Getenv("OKX_RFQ_ID"))
	if quoteId != "" {
		svc.QuoteId(quoteId)
	}
	if clQuoteId != "" {
		svc.ClQuoteId(clQuoteId)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cancel quote: quoteId=%s clQuoteId=%s sCode=%s sMsg=%s", ack.QuoteId, ack.ClQuoteId, ack.SCode, ack.SMsg)
}
