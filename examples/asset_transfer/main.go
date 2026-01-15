package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to asset transfer; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_AMT")
	from := os.Getenv("OKX_FROM")
	to := os.Getenv("OKX_TO")
	if ccy == "" || amt == "" || from == "" || to == "" {
		log.Fatal("missing env: OKX_CCY / OKX_AMT / OKX_FROM / OKX_TO")
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

	svc := c.NewAssetTransferService().Ccy(ccy).Amt(amt).From(from).To(to)
	if v := os.Getenv("OKX_TYPE"); v != "" {
		svc.Type(v)
	}
	if v := os.Getenv("OKX_SUB_ACCT"); v != "" {
		svc.SubAcct(v)
	}
	if v := os.Getenv("OKX_OMIT_POS_RISK"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("invalid env OKX_OMIT_POS_RISK=%q: %v", v, err)
		}
		svc.OmitPosRisk(b)
	}
	if v := os.Getenv("OKX_CLIENT_ID"); v != "" {
		svc.ClientId(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("asset_transfer: transId=%s clientId=%s ccy=%s amt=%s from=%s to=%s", ack.TransId, ack.ClientId, ack.Ccy, ack.Amt, ack.From, ack.To)
}
