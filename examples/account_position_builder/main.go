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

	svc := c.NewAccountPositionBuilderService()

	if acctLv := os.Getenv("OKX_ACCT_LV"); acctLv != "" {
		svc.AcctLv(acctLv)
	}
	if lever := os.Getenv("OKX_LEVER"); lever != "" {
		svc.Lever(lever)
	}
	if gt := os.Getenv("OKX_GREEKS_TYPE"); gt != "" {
		svc.GreeksType(gt)
	}
	if idxVol := os.Getenv("OKX_IDX_VOL"); idxVol != "" {
		svc.IdxVol(idxVol)
	}
	if incl := os.Getenv("OKX_INCL_REAL_POS_AND_EQ"); incl != "" {
		v, err := strconv.ParseBool(incl)
		if err != nil {
			log.Fatalf("invalid env OKX_INCL_REAL_POS_AND_EQ: %v", err)
		}
		svc.InclRealPosAndEq(v)
	}

	// 可选：注入一条模拟仓位（需要完整 3 个字段）。
	simPosInstId := os.Getenv("OKX_SIM_POS_INST_ID")
	simPosPos := os.Getenv("OKX_SIM_POS_POS")
	simPosAvgPx := os.Getenv("OKX_SIM_POS_AVG_PX")
	if simPosInstId != "" || simPosPos != "" || simPosAvgPx != "" {
		if simPosInstId == "" || simPosPos == "" || simPosAvgPx == "" {
			log.Fatal("missing env for simPos: OKX_SIM_POS_INST_ID / OKX_SIM_POS_POS / OKX_SIM_POS_AVG_PX")
		}
		p := okx.AccountPositionBuilderSimPos{
			InstId: simPosInstId,
			Pos:    simPosPos,
			AvgPx:  simPosAvgPx,
			Lever:  os.Getenv("OKX_SIM_POS_LEVER"), // optional
		}
		svc.SimPos([]okx.AccountPositionBuilderSimPos{p})
	}

	// 可选：注入一条模拟资产（需要完整 2 个字段）。
	simAssetCcy := os.Getenv("OKX_SIM_ASSET_CCY")
	simAssetAmt := os.Getenv("OKX_SIM_ASSET_AMT")
	if simAssetCcy != "" || simAssetAmt != "" {
		if simAssetCcy == "" || simAssetAmt == "" {
			log.Fatal("missing env for simAsset: OKX_SIM_ASSET_CCY / OKX_SIM_ASSET_AMT")
		}
		svc.SimAsset([]okx.AccountPositionBuilderSimAsset{{Ccy: simAssetCcy, Amt: simAssetAmt}})
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_position_builder: eq=%s marginRatio=%s totalImr=%s totalMmr=%s ts=%d", res.Eq, res.MarginRatio, res.TotalImr, res.TotalMmr, res.TS)
}
