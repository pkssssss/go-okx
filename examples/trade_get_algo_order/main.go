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

	algoId := os.Getenv("OKX_ALGO_ID")
	algoClOrdId := os.Getenv("OKX_ALGO_CL_ORD_ID")
	if algoId == "" && algoClOrdId == "" {
		log.Fatal("missing env: OKX_ALGO_ID or OKX_ALGO_CL_ORD_ID")
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

	svc := c.NewGetAlgoOrderService()
	if algoId != "" {
		svc.AlgoId(algoId)
	} else {
		svc.AlgoClOrdId(algoClOrdId)
	}

	order, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_get_algo_order: algoId=%s instId=%s ordType=%s side=%s state=%s cTime=%d uTime=%d", order.AlgoId, order.InstId, order.OrdType, order.Side, order.State, order.CTime, order.UTime)
}
