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

	svc := c.NewRFQRfqsService().
		RfqId(os.Getenv("OKX_RFQ_ID")).
		ClRfqId(os.Getenv("OKX_CL_RFQ_ID")).
		State(os.Getenv("OKX_RFQ_STATE")).
		BeginId(os.Getenv("OKX_BEGIN_ID")).
		EndId(os.Getenv("OKX_END_ID"))
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	rfqs, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rfqs=%d", len(rfqs))
	if len(rfqs) > 0 {
		log.Printf("first: rfqId=%s state=%s legs=%d", rfqs[0].RfqId, rfqs[0].State, len(rfqs[0].Legs))
	}
}
