package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	instId := os.Getenv("OKX_INST_ID")
	ordId := os.Getenv("OKX_ORD_ID")
	clOrdId := os.Getenv("OKX_CL_ORD_ID")

	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}
	if ordId == "" && clOrdId == "" {
		log.Fatal("missing env: OKX_ORD_ID or OKX_CL_ORD_ID")
	}

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

	ack, err := ws.CancelOrder(opCtx, okx.WSCancelOrderArg{
		InstId:  instId,
		OrdId:   ordId,
		ClOrdId: clOrdId,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ws cancel order ack: clOrdId=%s ordId=%s sCode=%s sMsg=%s ts=%d", ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg, ack.TS)
}
