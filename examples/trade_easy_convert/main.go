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
		log.Fatal("refusing to easy-convert; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	fromCcyListJSON := os.Getenv("OKX_FROM_CCY_LIST")
	toCcy := os.Getenv("OKX_TO_CCY")
	if fromCcyListJSON == "" || toCcy == "" {
		log.Fatal("missing env: OKX_FROM_CCY_LIST (JSON array) / OKX_TO_CCY")
	}
	source := os.Getenv("OKX_SOURCE") // optional: 1=交易账户, 2=资金账户

	var fromCcyList []string
	if err := json.Unmarshal([]byte(fromCcyListJSON), &fromCcyList); err != nil {
		log.Fatalf("invalid OKX_FROM_CCY_LIST: %v", err)
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

	svc := c.NewEasyConvertService().FromCcy(fromCcyList).ToCcy(toCcy)
	if source != "" {
		svc.Source(source)
	}

	acks, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_easy_convert: fromCcy=%v toCcy=%s source=%s acks=%d", fromCcyList, toCcy, source, len(acks))
	for i := 0; i < len(acks) && i < 5; i++ {
		x := acks[i]
		log.Printf("ack[%d]: fromCcy=%s fillFromSz=%s toCcy=%s fillToSz=%s status=%s uTime=%d", i, x.FromCcy, x.FillFromSz, x.ToCcy, x.FillToSz, x.Status, x.UTime)
	}
}
