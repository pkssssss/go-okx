package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to amend algo order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}

	algoId := os.Getenv("OKX_ALGO_ID")
	algoClOrdId := os.Getenv("OKX_ALGO_CL_ORD_ID")
	if algoId == "" && algoClOrdId == "" {
		log.Fatal("missing env: OKX_ALGO_ID or OKX_ALGO_CL_ORD_ID")
	}

	newSz := os.Getenv("OKX_NEW_SZ")
	if newSz == "" {
		log.Fatal("missing env: OKX_NEW_SZ")
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

	svc := c.NewAmendAlgoOrderService().InstId(instId).NewSz(newSz)
	if algoId != "" {
		svc.AlgoId(algoId)
	} else {
		svc.AlgoClOrdId(algoClOrdId)
	}
	if v := os.Getenv("OKX_REQ_ID"); v != "" {
		svc.ReqId(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_amend_algos: algoId=%s algoClOrdId=%s reqId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.ReqId, ack.SCode, ack.SMsg)
}
