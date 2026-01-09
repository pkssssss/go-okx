package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}
	if os.Getenv("OKX_CONFIRM_SET_COLLATERAL_ASSETS") != "YES" {
		log.Fatal("refusing to set collateral assets; set OKX_CONFIRM_SET_COLLATERAL_ASSETS=YES to continue")
	}

	typ := os.Getenv("OKX_TYPE") // all/custom
	if typ == "" {
		log.Fatal("missing env: OKX_TYPE")
	}

	enabledStr := os.Getenv("OKX_COLLATERAL_ENABLED")
	if enabledStr == "" {
		log.Fatal("missing env: OKX_COLLATERAL_ENABLED")
	}
	enabled, err := strconv.ParseBool(enabledStr)
	if err != nil {
		log.Fatalf("invalid env OKX_COLLATERAL_ENABLED: %v", err)
	}

	ccyListStr := os.Getenv("OKX_CCY_LIST") // optional, comma separated
	var ccyList []string
	if ccyListStr != "" {
		parts := strings.Split(ccyListStr, ",")
		for _, p := range parts {
			ccy := strings.TrimSpace(p)
			if ccy == "" {
				continue
			}
			ccyList = append(ccyList, ccy)
		}
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

	svc := c.NewAccountSetCollateralAssetsService().
		Type(typ).
		CollateralEnabled(enabled)
	if len(ccyList) > 0 {
		svc.CcyList(ccyList)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_set_collateral_assets: type=%s enabled=%v ccyList=%v", ack.Type, ack.CollateralEnabled, ack.CcyList)
}
