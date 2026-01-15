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
		log.Fatal("refusing to amend tradingBot signal tp/sl; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	tpSlType := os.Getenv("OKX_TP_SL_TYPE")
	if algoId == "" || tpSlType == "" {
		log.Fatal("missing env: OKX_ALGO_ID / OKX_TP_SL_TYPE")
	}

	exit := okx.TradingBotSignalExitSettingParam{
		TpSlType: tpSlType,
	}
	if v := os.Getenv("OKX_TP_PCT"); v != "" {
		exit.TpPct = v
	}
	if v := os.Getenv("OKX_SL_PCT"); v != "" {
		exit.SlPct = v
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

	ack, err := c.NewTradingBotSignalAmendTPSLService().AlgoId(algoId).ExitSettingParam(exit).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot signal amendTPSL: algoId=%s tpSlType=%s tpPct=%s slPct=%s", ack.AlgoId, exit.TpSlType, exit.TpPct, exit.SlPct)
}
