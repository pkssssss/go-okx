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
		log.Fatal("refusing to place order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	instId := os.Getenv("OKX_INST_ID")
	tdMode := os.Getenv("OKX_TD_MODE")
	side := os.Getenv("OKX_SIDE")
	ordType := os.Getenv("OKX_ORD_TYPE")
	sz := os.Getenv("OKX_SZ")
	px := os.Getenv("OKX_PX")
	pxUsd := os.Getenv("OKX_PX_USD")
	pxVol := os.Getenv("OKX_PX_VOL")

	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}
	if tdMode == "" {
		log.Fatal("missing env: OKX_TD_MODE")
	}
	if side == "" {
		log.Fatal("missing env: OKX_SIDE")
	}
	if ordType == "" {
		log.Fatal("missing env: OKX_ORD_TYPE")
	}
	if sz == "" {
		log.Fatal("missing env: OKX_SZ")
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

	ack, err := ws.PlaceOrder(opCtx, okx.WSPlaceOrderArg{
		InstId:  instId,
		TdMode:  tdMode,
		Side:    side,
		OrdType: ordType,
		Sz:      sz,
		Px:      px,
		PxUsd:   pxUsd,
		PxVol:   pxVol,
		ClOrdId: os.Getenv("OKX_CL_ORD_ID"),
		Tag:     os.Getenv("OKX_TAG"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ws place order ack: clOrdId=%s ordId=%s sCode=%s sMsg=%s ts=%d", ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg, ack.TS)
}
