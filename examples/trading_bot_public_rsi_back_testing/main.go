package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}
	timeframe := os.Getenv("OKX_TIMEFRAME")
	if timeframe == "" {
		timeframe = "1H"
	}
	thold := os.Getenv("OKX_THOLD")
	if thold == "" {
		thold = "70"
	}
	timePeriod := os.Getenv("OKX_TIME_PERIOD")
	if timePeriod == "" {
		timePeriod = "14"
	}

	svc := okx.NewClient().NewTradingBotPublicRSIBackTestingService().
		InstId(instId).
		Timeframe(timeframe).
		Thold(thold).
		TimePeriod(timePeriod)

	if v := os.Getenv("OKX_TRIGGER_COND"); v != "" {
		svc.TriggerCond(v)
	}
	if v := os.Getenv("OKX_DURATION"); v != "" {
		svc.Duration(v)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s timeframe=%s thold=%s timePeriod=%s triggerNum=%s\n", instId, timeframe, thold, timePeriod, res.TriggerNum)
}
