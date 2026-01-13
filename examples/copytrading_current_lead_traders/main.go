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

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
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

	items, err := c.NewCopyTradingCurrentLeadTradersService().InstType(instType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_current_lead_traders: instType=%s items=%d", instType, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	log.Printf("first: uniqueCode=%s nickName=%s margin=%s copyTotalAmt=%s copyTotalPnl=%s", it.UniqueCode, it.NickName, it.Margin, it.CopyTotalAmt, it.CopyTotalPnl)
}
