package main

import (
	"context"
	"log"
	"os"
	"strings"

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
		log.Fatal("refusing to cancel batch quotes; set OKX_CONFIRM=YES to continue")
	}

	quoteIdsStr := os.Getenv("OKX_QUOTE_IDS")      // optional, comma separated
	clQuoteIdsStr := os.Getenv("OKX_CL_QUOTE_IDS") // optional, comma separated
	if quoteIdsStr == "" && clQuoteIdsStr == "" {
		log.Fatal("missing env: OKX_QUOTE_IDS or OKX_CL_QUOTE_IDS")
	}

	quoteIds := splitCommaList(quoteIdsStr)
	clQuoteIds := splitCommaList(clQuoteIdsStr)

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

	svc := c.NewRFQCancelBatchQuotesService()
	if len(quoteIds) > 0 {
		svc.QuoteIds(quoteIds)
	}
	if len(clQuoteIds) > 0 {
		svc.ClQuoteIds(clQuoteIds)
	}

	acks, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rfq cancel-batch-quotes: acks=%d", len(acks))
	for i, ack := range acks {
		log.Printf("ack[%d]: quoteId=%s clQuoteId=%s sCode=%s sMsg=%s", i, ack.QuoteId, ack.ClQuoteId, ack.SCode, ack.SMsg)
	}
}

func splitCommaList(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		v := strings.TrimSpace(p)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}
