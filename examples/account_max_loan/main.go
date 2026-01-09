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

	demo := os.Getenv("OKX_DEMO") == "1"

	mgnMode := os.Getenv("OKX_MGN_MODE")
	if mgnMode == "" {
		mgnMode = "cross"
	}

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	ccy := os.Getenv("OKX_CCY")
	mgnCcy := os.Getenv("OKX_MGN_CCY")
	tradeQuoteCcy := os.Getenv("OKX_TRADE_QUOTE_CCY")

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

	svc := c.NewAccountMaxLoanService().MgnMode(mgnMode).InstId(instId)
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if mgnCcy != "" {
		svc.MgnCcy(mgnCcy)
	}
	if tradeQuoteCcy != "" {
		svc.TradeQuoteCcy(tradeQuoteCcy)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_max_loan: items=%d", len(data))
	for _, it := range data {
		log.Printf("instId=%s mgnMode=%s mgnCcy=%s ccy=%s side=%s maxLoan=%s", it.InstId, it.MgnMode, it.MgnCcy, it.Ccy, it.Side, it.MaxLoan)
	}
}
