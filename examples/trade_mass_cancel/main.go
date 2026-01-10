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
		log.Fatal("refusing to mass-cancel MMP orders; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		log.Fatal("missing env: OKX_INST_FAMILY (e.g. BTC-USD)")
	}
	lockInterval := os.Getenv("OKX_LOCK_INTERVAL") // optional: 0-10000 ms

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

	svc := c.NewMassCancelService().InstType("OPTION").InstFamily(instFamily)
	if lockInterval != "" {
		svc.LockInterval(lockInterval)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("trade_mass_cancel: instType=OPTION instFamily=%s lockInterval=%s result=%v", instFamily, lockInterval, ack.Result)
}
