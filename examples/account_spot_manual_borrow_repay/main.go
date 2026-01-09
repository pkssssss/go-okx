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
	if os.Getenv("OKX_CONFIRM_SPOT_MANUAL_BORROW_REPAY") != "YES" {
		log.Fatal("refusing to manual borrow/repay; set OKX_CONFIRM_SPOT_MANUAL_BORROW_REPAY=YES to continue")
	}

	ccy := os.Getenv("OKX_CCY")
	side := os.Getenv("OKX_SIDE") // borrow/repay
	amt := os.Getenv("OKX_AMT")
	if ccy == "" || side == "" || amt == "" {
		log.Fatal("missing env: OKX_CCY / OKX_SIDE / OKX_AMT")
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

	ack, err := c.NewAccountSpotManualBorrowRepayService().
		Ccy(ccy).
		Side(side).
		Amt(amt).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("spot_manual_borrow_repay: ccy=%s side=%s amt=%s", ack.Ccy, ack.Side, ack.Amt)
}
