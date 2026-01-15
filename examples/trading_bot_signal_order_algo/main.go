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
		log.Fatal("refusing to create tradingBot signal strategy; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	signalChanId := os.Getenv("OKX_SIGNAL_CHAN_ID")
	lever := os.Getenv("OKX_LEVER")
	investAmt := os.Getenv("OKX_INVEST_AMT")
	subOrdType := os.Getenv("OKX_SUB_ORD_TYPE")
	if signalChanId == "" || lever == "" || investAmt == "" || subOrdType == "" {
		log.Fatal("missing env: OKX_SIGNAL_CHAN_ID / OKX_LEVER / OKX_INVEST_AMT / OKX_SUB_ORD_TYPE")
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

	svc := c.NewTradingBotSignalOrderAlgoService().
		SignalChanId(signalChanId).
		Lever(lever).
		InvestAmt(investAmt).
		SubOrdType(subOrdType)

	includeAllRaw := os.Getenv("OKX_INCLUDE_ALL")
	includeAll := false
	if includeAllRaw != "" {
		b, err := strconv.ParseBool(includeAllRaw)
		if err != nil {
			log.Fatalf("invalid env OKX_INCLUDE_ALL=%q: %v", includeAllRaw, err)
		}
		includeAll = b
		svc.IncludeAll(b)
	}

	if ratio := os.Getenv("OKX_RATIO"); ratio != "" {
		svc.Ratio(ratio)
	}

	if !includeAll {
		instIdsStr := os.Getenv("OKX_INST_IDS")
		if instIdsStr == "" {
			log.Fatal("missing env: OKX_INST_IDS (comma separated, required when OKX_INCLUDE_ALL is not true)")
		}

		var instIds []string
		for _, p := range strings.Split(instIdsStr, ",") {
			id := strings.TrimSpace(p)
			if id == "" {
				continue
			}
			instIds = append(instIds, id)
		}
		if len(instIds) == 0 {
			log.Fatal("missing env: OKX_INST_IDS (comma separated, required when OKX_INCLUDE_ALL is not true)")
		}
		svc.InstIds(instIds)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot signal order-algo: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
