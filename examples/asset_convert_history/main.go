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

	demo := os.Getenv("OKX_DEMO") == "1"

	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")
	tag := os.Getenv("OKX_TAG")
	clTReqId := os.Getenv("OKX_CL_T_REQ_ID")

	var limit *int
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		limit = &n
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

	svc := c.NewAssetConvertHistoryService()
	if clTReqId != "" {
		svc.ClTReqId(clTReqId)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}
	if limit != nil {
		svc.Limit(*limit)
	}
	if tag != "" {
		svc.Tag(tag)
	}

	trades, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("convert_history=%d", len(trades))
	for i := 0; i < len(trades) && i < 5; i++ {
		x := trades[i]
		log.Printf("trade[%d]: tradeId=%s instId=%s side=%s state=%s fillPx=%s fillBaseSz=%s fillQuoteSz=%s ts=%d", i, x.TradeId, x.InstId, x.Side, x.State, x.FillPx, x.FillBaseSz, x.FillQuoteSz, x.TS)
	}
}
