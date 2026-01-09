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

	instType := os.Getenv("OKX_INST_TYPE") // MARGIN / SWAP / FUTURES
	mgnMode := os.Getenv("OKX_MGN_MODE")   // isolated / cross
	lever := os.Getenv("OKX_LEVER")
	if instType == "" || mgnMode == "" || lever == "" {
		log.Fatal("missing env: OKX_INST_TYPE / OKX_MGN_MODE / OKX_LEVER")
	}

	instId := os.Getenv("OKX_INST_ID")   // optional
	ccy := os.Getenv("OKX_CCY")          // optional
	posSide := os.Getenv("OKX_POS_SIDE") // optional
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

	svc := c.NewAccountAdjustLeverageInfoService().
		InstType(instType).
		MgnMode(mgnMode).
		Lever(lever)
	if instId != "" {
		svc.InstId(instId)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if posSide != "" {
		svc.PosSide(posSide)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_adjust_leverage_info: existOrd=%v maxLever=%s minLever=%s estMgn=%s estLiqPx=%s estMaxAmt=%s",
		res.ExistOrd, res.MaxLever, res.MinLever, res.EstMgn, res.EstLiqPx, res.EstMaxAmt,
	)
}
