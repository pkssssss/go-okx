package main

import (
	"context"
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
		log.Fatal("refusing to amend sprd order via ws; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	ordId := os.Getenv("OKX_ORD_ID")
	clOrdId := os.Getenv("OKX_CL_ORD_ID")
	if ordId == "" && clOrdId == "" {
		log.Fatal("missing env: OKX_ORD_ID or OKX_CL_ORD_ID")
	}

	newSz := os.Getenv("OKX_NEW_SZ")
	newPx := os.Getenv("OKX_NEW_PX")
	if newSz == "" && newPx == "" {
		log.Fatal("missing env: OKX_NEW_SZ or OKX_NEW_PX")
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

	ws := c.NewWSBusinessPrivate()

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

	ack, err := ws.SprdAmendOrder(opCtx, okx.WSSprdAmendOrderArg{
		OrdId:   ordId,
		ClOrdId: clOrdId,
		ReqId:   os.Getenv("OKX_REQ_ID"),
		NewSz:   newSz,
		NewPx:   newPx,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ws sprd amend order ack: reqId=%s clOrdId=%s ordId=%s sCode=%s sMsg=%s ts=%d", ack.ReqId, ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg, ack.TS)
}
