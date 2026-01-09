package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel algo orders; set OKX_CONFIRM=YES to continue")
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

	acks, err := c.NewCancelAlgoOrdersService().
		Orders([]okx.CancelAlgoOrder{{InstId: instId, AlgoId: algoId, AlgoClOrdId: algoClOrdId}}).
		Do(context.Background())
	if err != nil {
		var batchErr *okx.TradeAlgoBatchError
		if errors.As(err, &batchErr) {
			log.Printf("cancel algo orders partial failure: %v", err)
			for i, ack := range batchErr.Acks {
				log.Printf("ack[%d]: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", i, ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
			}
			return
		}
		log.Fatal(err)
	}

	log.Printf("trade_cancel_algos: acks=%d", len(acks))
	for i, ack := range acks {
		log.Printf("ack[%d]: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", i, ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
	}
}
