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
	if os.Getenv("OKX_CONFIRM_SET_LEVERAGE") != "YES" {
		log.Fatal("refusing to set leverage; set OKX_CONFIRM_SET_LEVERAGE=YES to continue")
	}

	lever := os.Getenv("OKX_LEVER")
	mgnMode := os.Getenv("OKX_MGN_MODE")
	if lever == "" || mgnMode == "" {
		log.Fatal("missing env: OKX_LEVER / OKX_MGN_MODE")
	}

	instId := os.Getenv("OKX_INST_ID")
	ccy := os.Getenv("OKX_CCY")
	if instId == "" && ccy == "" {
		log.Fatal("missing env: OKX_INST_ID or OKX_CCY")
	}
	if instId != "" && ccy != "" {
		log.Fatal("invalid env: OKX_INST_ID and OKX_CCY are mutually exclusive")
	}

	posSide := os.Getenv("OKX_POS_SIDE") // optional: long/short
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

	svc := c.NewAccountSetLeverageService().
		Lever(lever).
		MgnMode(mgnMode)
	if instId != "" {
		svc.InstId(instId)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if posSide != "" {
		svc.PosSide(posSide)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("set_leverage: instId=%s mgnMode=%s posSide=%s lever=%s", ack.InstId, ack.MgnMode, ack.PosSide, ack.Lever)
}
