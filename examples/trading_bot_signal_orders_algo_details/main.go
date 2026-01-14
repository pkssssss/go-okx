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

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	if algoId == "" {
		log.Fatal("missing env: OKX_ALGO_ID")
	}

	algoOrdType := os.Getenv("OKX_ALGO_ORD_TYPE")
	if algoOrdType == "" {
		algoOrdType = "contract"
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

	o, err := c.NewTradingBotSignalOrdersAlgoDetailsService().AlgoOrdType(algoOrdType).AlgoId(algoId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("signal order: algoId=%s algoOrdType=%s state=%s instType=%s instIds=%d lever=%s investAmt=%s subOrdType=%s signalChanId=%s cTime=%d uTime=%d",
		o.AlgoId, o.AlgoOrdType, o.State, o.InstType, len(o.InstIds), o.Lever, o.InvestAmt, o.SubOrdType, o.SignalChanId, o.CTime, o.UTime,
	)
}
