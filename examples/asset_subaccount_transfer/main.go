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
	if os.Getenv("OKX_CONFIRM_SUBACCOUNT_TRANSFER") != "YES" {
		log.Fatal("refusing to subaccount transfer; set OKX_CONFIRM_SUBACCOUNT_TRANSFER=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_AMT")
	from := os.Getenv("OKX_FROM")
	to := os.Getenv("OKX_TO")
	fromSub := os.Getenv("OKX_FROM_SUBACCOUNT")
	toSub := os.Getenv("OKX_TO_SUBACCOUNT")
	if ccy == "" || amt == "" || from == "" || to == "" || fromSub == "" || toSub == "" {
		log.Fatal("missing env: OKX_CCY / OKX_AMT / OKX_FROM / OKX_TO / OKX_FROM_SUBACCOUNT / OKX_TO_SUBACCOUNT")
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

	ack, err := c.NewAssetSubaccountTransferService().
		Ccy(ccy).
		Amt(amt).
		From(from).
		To(to).
		FromSubAccount(fromSub).
		ToSubAccount(toSub).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("subaccount_transfer: transId=%s", ack.TransId)
}
