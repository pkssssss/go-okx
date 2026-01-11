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
		log.Fatal("refusing to create rfq; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	counterpartiesJSON := os.Getenv("OKX_RFQ_COUNTERPARTIES") // JSON array
	legsJSON := os.Getenv("OKX_RFQ_LEGS")                     // JSON array
	if counterpartiesJSON == "" || legsJSON == "" {
		log.Fatal("missing env: OKX_RFQ_COUNTERPARTIES / OKX_RFQ_LEGS (JSON arrays)")
	}

	var counterparties []string
	if err := json.Unmarshal([]byte(counterpartiesJSON), &counterparties); err != nil {
		log.Fatalf("invalid OKX_RFQ_COUNTERPARTIES: %v", err)
	}

	var legs []okx.RFQLeg
	if err := json.Unmarshal([]byte(legsJSON), &legs); err != nil {
		log.Fatalf("invalid OKX_RFQ_LEGS: %v", err)
	}

	var acctAlloc []okx.RFQAcctAlloc
	if v := os.Getenv("OKX_RFQ_ACCT_ALLOC"); v != "" {
		if err := json.Unmarshal([]byte(v), &acctAlloc); err != nil {
			log.Fatalf("invalid OKX_RFQ_ACCT_ALLOC: %v", err)
		}
	}

	anonymous, hasAnonymous := os.LookupEnv("OKX_RFQ_ANONYMOUS") // "true"/"false"
	allowPartial, hasAllowPartial := os.LookupEnv("OKX_RFQ_ALLOW_PARTIAL_EXECUTION")

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

	svc := c.NewRFQCreateRFQService().
		Counterparties(counterparties).
		Legs(legs).
		ClRfqId(os.Getenv("OKX_RFQ_CL_RFQ_ID")).
		Tag(os.Getenv("OKX_RFQ_TAG"))
	if hasAnonymous {
		svc.Anonymous(anonymous == "true" || anonymous == "1")
	}
	if hasAllowPartial {
		svc.AllowPartialExecution(allowPartial == "true" || allowPartial == "1")
	}
	if len(acctAlloc) > 0 {
		svc.AcctAlloc(acctAlloc)
	}

	rfq, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("create rfq: rfqId=%s state=%s cTime=%d legs=%d", rfq.RfqId, rfq.State, rfq.CTime, len(rfq.Legs))
}
