package main

import (
	"context"
	"fmt"
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

	svc := c.NewTradingBotSignalOrdersAlgoPendingService().AlgoOrdType(algoOrdType)
	if v := os.Getenv("OKX_ALGO_ID"); v != "" {
		svc.AlgoId(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("algoOrdType=%s count=%d\n", algoOrdType, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	fmt.Printf("first algoId=%s state=%s instType=%s instIds=%d signalChanId=%s cTime=%d uTime=%d\n",
		it.AlgoId, it.State, it.InstType, len(it.InstIds), it.SignalChanId, it.CTime, it.UTime,
	)
}
