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
		log.Fatal("refusing to cancel rfq; set OKX_CONFIRM=YES to continue")
	}

	rfqId := os.Getenv("OKX_RFQ_ID")
	clRfqId := os.Getenv("OKX_CL_RFQ_ID")
	if rfqId == "" && clRfqId == "" {
		log.Fatal("missing env: OKX_RFQ_ID or OKX_CL_RFQ_ID")
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

	svc := c.NewRFQCancelRFQService()
	if rfqId != "" {
		svc.RfqId(rfqId)
	}
	if clRfqId != "" {
		svc.ClRfqId(clRfqId)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("cancel rfq: rfqId=%s clRfqId=%s sCode=%s sMsg=%s", ack.RfqId, ack.ClRfqId, ack.SCode, ack.SMsg)
}
