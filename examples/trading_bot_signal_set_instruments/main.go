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

	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set tradingBot signal instruments; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	algoId := os.Getenv("OKX_ALGO_ID")
	if algoId == "" {
		log.Fatal("missing env: OKX_ALGO_ID")
	}

	includeAllRaw := os.Getenv("OKX_INCLUDE_ALL")
	if includeAllRaw == "" {
		log.Fatal("missing env: OKX_INCLUDE_ALL (true/false)")
	}
	includeAll, err := strconv.ParseBool(includeAllRaw)
	if err != nil {
		log.Fatalf("invalid env OKX_INCLUDE_ALL=%q: %v", includeAllRaw, err)
	}

	instIdsStr := os.Getenv("OKX_INST_IDS")
	var instIds []string
	if !includeAll {
		if instIdsStr == "" {
			log.Fatal("missing env: OKX_INST_IDS (comma separated, required when OKX_INCLUDE_ALL=false)")
		}
		for _, p := range strings.Split(instIdsStr, ",") {
			id := strings.TrimSpace(p)
			if id == "" {
				continue
			}
			instIds = append(instIds, id)
		}
		if len(instIds) == 0 {
			log.Fatal("missing env: OKX_INST_IDS (comma separated, required when OKX_INCLUDE_ALL=false)")
		}
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

	svc := c.NewTradingBotSignalSetInstrumentsService().AlgoId(algoId).IncludeAll(includeAll)
	if len(instIds) > 0 {
		svc.InstIds(instIds)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot signal set-instruments: algoId=%s includeAll=%t instIds=%v", ack.AlgoId, includeAll, instIds)
}
