package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}
	if os.Getenv("OKX_CONFIRM_SET_AUTO_LOAN") != "YES" {
		log.Fatal("refusing to set auto loan; set OKX_CONFIRM_SET_AUTO_LOAN=YES to continue")
	}

	autoLoanStr := os.Getenv("OKX_AUTO_LOAN")
	if autoLoanStr == "" {
		log.Fatal("missing env: OKX_AUTO_LOAN")
	}
	autoLoan, err := strconv.ParseBool(autoLoanStr)
	if err != nil {
		log.Fatalf("invalid env OKX_AUTO_LOAN: %v", err)
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

	ack, err := c.NewAccountSetAutoLoanService().AutoLoan(autoLoan).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_auto_loan: autoLoan=%v", ack.AutoLoan)
}
