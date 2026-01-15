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
		log.Fatal("refusing to stop tradingBot grid orders; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	instId := os.Getenv("OKX_INST_ID")
	stopType := os.Getenv("OKX_STOP_TYPE")
	if algoId == "" || instId == "" || stopType == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_INST_ID / OKX_STOP_TYPE")
	}

	algoOrdType := os.Getenv("OKX_ALGO_ORD_TYPE")
	if algoOrdType == "" {
		algoOrdType = "grid"
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

	acks, err := c.NewTradingBotGridStopOrderAlgoService().
		Orders([]okx.TradingBotGridStopOrder{{
			AlgoId:      algoId,
			InstId:      instId,
			AlgoOrdType: algoOrdType,
			StopType:    stopType,
		}}).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot grid stop-order-algo: acks=%d", len(acks))
	if len(acks) > 0 {
		ack := acks[0]
		log.Printf("first: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
	}
}
