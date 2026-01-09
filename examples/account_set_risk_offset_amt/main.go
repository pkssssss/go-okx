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
	if os.Getenv("OKX_CONFIRM_SET_RISK_OFFSET_AMT") != "YES" {
		log.Fatal("refusing to set risk offset amt; set OKX_CONFIRM_SET_RISK_OFFSET_AMT=YES to continue")
	}

	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_CL_SPOT_IN_USE_AMT")
	if ccy == "" || amt == "" {
		log.Fatal("missing env: OKX_CCY / OKX_CL_SPOT_IN_USE_AMT")
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

	ack, err := c.NewAccountSetRiskOffsetAmtService().Ccy(ccy).ClSpotInUseAmt(amt).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_risk_offset_amt: ccy=%s clSpotInUseAmt=%s", ack.Ccy, ack.ClSpotInUseAmt)
}
