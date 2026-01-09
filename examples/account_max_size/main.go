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

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	tdMode := os.Getenv("OKX_TD_MODE")
	if tdMode == "" {
		tdMode = "cash"
	}

	ccy := os.Getenv("OKX_CCY")
	px := os.Getenv("OKX_PX")
	leverage := os.Getenv("OKX_LEVERAGE")
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

	svc := c.NewAccountMaxSizeService().
		InstId(instId).
		TdMode(tdMode)
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if px != "" {
		svc.Px(px)
	}
	if leverage != "" {
		svc.Leverage(leverage)
	}
	if tradeQuoteCcy != "" {
		svc.TradeQuoteCcy(tradeQuoteCcy)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_max_size: items=%d", len(data))
	for _, it := range data {
		log.Printf("instId=%s ccy=%s maxBuy=%s maxSell=%s", it.InstId, it.Ccy, it.MaxBuy, it.MaxSell)
	}
}
