package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set first copy settings; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	uniqueCode := os.Getenv("OKX_UNIQUE_CODE")
	copyMgnMode := os.Getenv("OKX_COPY_MGN_MODE")
	copyInstIdType := os.Getenv("OKX_COPY_INST_ID_TYPE")
	copyTotalAmt := os.Getenv("OKX_COPY_TOTAL_AMT")
	subPosCloseType := os.Getenv("OKX_SUB_POS_CLOSE_TYPE")
	if uniqueCode == "" || copyMgnMode == "" || copyInstIdType == "" || copyTotalAmt == "" || subPosCloseType == "" {
		log.Fatal("missing env: OKX_UNIQUE_CODE / OKX_COPY_MGN_MODE / OKX_COPY_INST_ID_TYPE / OKX_COPY_TOTAL_AMT / OKX_SUB_POS_CLOSE_TYPE")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	copyMode := os.Getenv("OKX_COPY_MODE")
	copyAmt := os.Getenv("OKX_COPY_AMT")
	copyRatio := os.Getenv("OKX_COPY_RATIO")
	if copyMode == "" || copyMode == "fixed_amount" {
		if copyAmt == "" {
			log.Fatal("missing env: OKX_COPY_AMT (required for fixed_amount/default)")
		}
	} else if copyMode == "ratio_copy" {
		if copyRatio == "" {
			log.Fatal("missing env: OKX_COPY_RATIO (required for ratio_copy)")
		}
	}

	instId := os.Getenv("OKX_INST_ID")

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

	svc := c.NewCopyTradingFirstCopySettingsService().
		InstType(instType).
		UniqueCode(uniqueCode).
		CopyMgnMode(copyMgnMode).
		CopyInstIdType(copyInstIdType).
		CopyTotalAmt(copyTotalAmt).
		SubPosCloseType(subPosCloseType)

	if instId != "" {
		svc.InstId(instId)
	}
	if copyMode != "" {
		svc.CopyMode(copyMode)
	}
	if copyAmt != "" {
		svc.CopyAmt(copyAmt)
	}
	if copyRatio != "" {
		svc.CopyRatio(copyRatio)
	}
	if v := os.Getenv("OKX_TP_RATIO"); v != "" {
		svc.TpRatio(v)
	}
	if v := os.Getenv("OKX_SL_RATIO"); v != "" {
		svc.SlRatio(v)
	}
	if v := os.Getenv("OKX_SL_TOTAL_AMT"); v != "" {
		svc.SlTotalAmt(v)
	}
	if v := os.Getenv("OKX_TAG"); v != "" {
		svc.Tag(v)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_first_copy_settings: result=%v", res.Result)
}
