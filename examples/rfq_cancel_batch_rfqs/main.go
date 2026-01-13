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
		log.Fatal("refusing to cancel batch rfqs; set OKX_CONFIRM=YES to continue")
	}

	rfqIdsStr := os.Getenv("OKX_RFQ_IDS")      // optional, comma separated
	clRfqIdsStr := os.Getenv("OKX_CL_RFQ_IDS") // optional, comma separated
	if rfqIdsStr == "" && clRfqIdsStr == "" {
		log.Fatal("missing env: OKX_RFQ_IDS or OKX_CL_RFQ_IDS")
	}

	rfqIds := splitCommaList(rfqIdsStr)
	clRfqIds := splitCommaList(clRfqIdsStr)

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

	svc := c.NewRFQCancelBatchRFQsService()
	if len(rfqIds) > 0 {
		svc.RfqIds(rfqIds)
	}
	if len(clRfqIds) > 0 {
		svc.ClRfqIds(clRfqIds)
	}

	acks, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rfq cancel-batch-rfqs: acks=%d", len(acks))
	for i, ack := range acks {
		log.Printf("ack[%d]: rfqId=%s clRfqId=%s sCode=%s sMsg=%s", i, ack.RfqId, ack.ClRfqId, ack.SCode, ack.SMsg)
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
