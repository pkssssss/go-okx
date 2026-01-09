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

	instType := os.Getenv("OKX_INST_TYPE") // required: SPOT/MARGIN/SWAP/FUTURES/OPTION
	if instType == "" {
		log.Fatal("missing env: OKX_INST_TYPE")
	}

	instFamily := os.Getenv("OKX_INST_FAMILY") // optional (OPTION required)
	instId := os.Getenv("OKX_INST_ID")         // optional
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

	svc := c.NewAccountInstrumentsService().InstType(instType)
	if instFamily != "" {
		svc.InstFamily(instFamily)
	}
	if instId != "" {
		svc.InstId(instId)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_instruments: items=%d", len(data))
	if len(data) > 0 {
		it := data[0]
		log.Printf("first: instType=%s instId=%s instFamily=%s tickSz=%s lotSz=%s minSz=%s ctVal=%s ctMult=%s groupId=%s",
			it.InstType, it.InstId, it.InstFamily, it.TickSz, it.LotSz, it.MinSz, it.CtVal, it.CtMult, it.GroupId,
		)
	}
}
