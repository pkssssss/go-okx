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

	mgnMode := os.Getenv("OKX_MGN_MODE")
	if mgnMode == "" {
		mgnMode = "cross"
	}
	instId := os.Getenv("OKX_INST_ID")
	ccy := os.Getenv("OKX_CCY")
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

	svc := c.NewAccountLeverageInfoService().MgnMode(mgnMode)
	if instId != "" {
		svc.InstId(instId)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("leverage_info=%d mgnMode=%s instId=%s ccy=%s", len(data), mgnMode, instId, ccy)
	for i := 0; i < len(data) && i < 10; i++ {
		x := data[i]
		log.Printf("item[%d]: instId=%s ccy=%s mgnMode=%s posSide=%s lever=%s", i, x.InstId, x.Ccy, x.MgnMode, x.PosSide, x.Lever)
	}
}
