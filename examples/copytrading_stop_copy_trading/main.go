package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to stop copy trading; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	uniqueCode := os.Getenv("OKX_UNIQUE_CODE")
	if uniqueCode == "" {
		log.Fatal("missing env: OKX_UNIQUE_CODE")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	subPosCloseType := os.Getenv("OKX_SUB_POS_CLOSE_TYPE")

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

	svc := c.NewCopyTradingStopCopyTradingService().
		InstType(instType).
		UniqueCode(uniqueCode)
	if subPosCloseType != "" {
		svc.SubPosCloseType(subPosCloseType)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_stop_copy_trading: result=%v", res.Result)
}
