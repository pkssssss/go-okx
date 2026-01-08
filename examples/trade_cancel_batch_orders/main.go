package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/pkssssss/go-okx"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	ordersJSON := os.Getenv("OKX_BATCH_CANCEL_ORDERS")
	if ordersJSON == "" {
		log.Fatal("missing env: OKX_BATCH_CANCEL_ORDERS (JSON array)")
	}
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel orders; set OKX_CONFIRM=YES to continue")
	}

	var orders []okx.BatchCancelOrder
	if err := json.Unmarshal([]byte(ordersJSON), &orders); err != nil {
		log.Fatalf("invalid OKX_BATCH_CANCEL_ORDERS: %v", err)
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

	acks, err := c.NewBatchCancelOrdersService().Orders(orders).Do(context.Background())
	if err != nil {
		var batchErr *okx.TradeBatchError
		if errors.As(err, &batchErr) {
			log.Printf("batch cancel orders partial failure: %v", err)
			for i, ack := range batchErr.Acks {
				log.Printf("ack[%d]: clOrdId=%s ordId=%s sCode=%s sMsg=%s", i, ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg)
			}
			return
		}
		log.Fatal(err)
	}

	log.Printf("batch cancel orders acks=%d", len(acks))
	for i, ack := range acks {
		log.Printf("ack[%d]: clOrdId=%s ordId=%s sCode=%s sMsg=%s", i, ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg)
	}
}
