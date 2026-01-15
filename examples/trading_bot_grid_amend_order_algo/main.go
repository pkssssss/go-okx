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

	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to amend tradingBot grid order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	instId := os.Getenv("OKX_INST_ID")
	if algoId == "" || instId == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_INST_ID")
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

	svc := c.NewTradingBotGridAmendOrderAlgoService().AlgoId(algoId).InstId(instId)

	hasUpdate := false
	if v, ok := os.LookupEnv("OKX_SL_TRIGGER_PX"); ok {
		svc.SlTriggerPx(v)
		hasUpdate = true
	}
	if v, ok := os.LookupEnv("OKX_TP_TRIGGER_PX"); ok {
		svc.TpTriggerPx(v)
		hasUpdate = true
	}
	if v, ok := os.LookupEnv("OKX_TP_RATIO"); ok {
		svc.TpRatio(v)
		hasUpdate = true
	}
	if v, ok := os.LookupEnv("OKX_SL_RATIO"); ok {
		svc.SlRatio(v)
		hasUpdate = true
	}
	if v := os.Getenv("OKX_TOP_UP_AMT"); v != "" {
		svc.TopUpAmt(v)
		hasUpdate = true
	}

	if !hasUpdate {
		log.Fatal("missing update env: set one of OKX_SL_TRIGGER_PX / OKX_TP_TRIGGER_PX / OKX_TP_RATIO / OKX_SL_RATIO / OKX_TOP_UP_AMT")
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid amend-order-algo: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
