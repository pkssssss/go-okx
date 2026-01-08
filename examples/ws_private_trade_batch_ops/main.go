package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

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
		log.Fatal("refusing to send batch trade ops; set OKX_CONFIRM=YES to continue")
	}

	op := os.Getenv("OKX_WS_BATCH_OP")
	if op == "" {
		log.Fatal("missing env: OKX_WS_BATCH_OP (order/cancel-order/amend-order)")
	}
	argsJSON := os.Getenv("OKX_WS_BATCH_ARGS")
	if argsJSON == "" {
		log.Fatal("missing env: OKX_WS_BATCH_ARGS (JSON array)")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	timeout := 10 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
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

	ws := c.NewWSPrivate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}
	defer func() {
		ws.Close()
		<-ws.Done()
	}()

	opCtx, opCancel := context.WithTimeout(context.Background(), timeout)
	defer opCancel()

	switch op {
	case "order":
		var args []okx.WSPlaceOrderArg
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			log.Fatalf("invalid OKX_WS_BATCH_ARGS: %v", err)
		}
		acks, err := ws.PlaceOrders(opCtx, args...)
		if err != nil {
			var batchErr *okx.WSTradeOpBatchError
			if errors.As(err, &batchErr) {
				log.Printf("ws batch op partial failure: %v", err)
				printAcks(batchErr.Acks)
				return
			}
			log.Fatal(err)
		}
		printAcks(acks)
	case "cancel-order":
		var args []okx.WSCancelOrderArg
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			log.Fatalf("invalid OKX_WS_BATCH_ARGS: %v", err)
		}
		acks, err := ws.CancelOrders(opCtx, args...)
		if err != nil {
			var batchErr *okx.WSTradeOpBatchError
			if errors.As(err, &batchErr) {
				log.Printf("ws batch op partial failure: %v", err)
				printAcks(batchErr.Acks)
				return
			}
			log.Fatal(err)
		}
		printAcks(acks)
	case "amend-order":
		var args []okx.WSAmendOrderArg
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			log.Fatalf("invalid OKX_WS_BATCH_ARGS: %v", err)
		}
		acks, err := ws.AmendOrders(opCtx, args...)
		if err != nil {
			var batchErr *okx.WSTradeOpBatchError
			if errors.As(err, &batchErr) {
				log.Printf("ws batch op partial failure: %v", err)
				printAcks(batchErr.Acks)
				return
			}
			log.Fatal(err)
		}
		printAcks(acks)
	default:
		log.Fatalf("invalid OKX_WS_BATCH_OP: %s (order/cancel-order/amend-order)", op)
	}
}

func printAcks(acks []okx.TradeOrderAck) {
	log.Printf("acks=%d", len(acks))
	for i, ack := range acks {
		log.Printf("ack[%d]: clOrdId=%s ordId=%s reqId=%s sCode=%s sMsg=%s ts=%d", i, ack.ClOrdId, ack.OrdId, ack.ReqId, ack.SCode, ack.SMsg, ack.TS)
	}
}
