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
	if os.Getenv("OKX_CONFIRM_POSITION_MARGIN_BALANCE") != "YES" {
		log.Fatal("refusing to adjust margin; set OKX_CONFIRM_POSITION_MARGIN_BALANCE=YES to continue")
	}

	instId := os.Getenv("OKX_INST_ID")
	posSide := os.Getenv("OKX_POS_SIDE")
	typ := os.Getenv("OKX_TYPE") // add/reduce
	amt := os.Getenv("OKX_AMT")
	if instId == "" || posSide == "" || typ == "" || amt == "" {
		log.Fatal("missing env: OKX_INST_ID / OKX_POS_SIDE / OKX_TYPE / OKX_AMT")
	}

	ccy := os.Getenv("OKX_CCY") // optional
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

	svc := c.NewAccountPositionMarginBalanceService().
		InstId(instId).
		PosSide(posSide).
		Type(typ).
		Amt(amt)
	if ccy != "" {
		svc.Ccy(ccy)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("position_margin_balance: instId=%s posSide=%s type=%s amt=%s leverage=%s ccy=%s", ack.InstId, ack.PosSide, ack.Type, ack.Amt, ack.Leverage, ack.Ccy)
}
