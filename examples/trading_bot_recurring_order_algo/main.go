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

	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to place tradingBot recurring order; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	stgyName := os.Getenv("OKX_STGY_NAME")
	if stgyName == "" {
		stgyName = "my_recurring"
	}
	period := os.Getenv("OKX_PERIOD")
	if period == "" {
		period = "daily"
	}
	recurringTime := os.Getenv("OKX_RECURRING_TIME")
	if recurringTime == "" {
		recurringTime = "00:00"
	}
	timeZone := os.Getenv("OKX_TIME_ZONE")
	if timeZone == "" {
		timeZone = "UTC"
	}
	amt := os.Getenv("OKX_AMT")
	if amt == "" {
		log.Fatal("missing env: OKX_AMT")
	}
	investmentCcy := os.Getenv("OKX_INVESTMENT_CCY")
	if investmentCcy == "" {
		investmentCcy = "USDT"
	}
	tdMode := os.Getenv("OKX_TD_MODE")
	if tdMode == "" {
		tdMode = "cash"
	}

	listCcy := os.Getenv("OKX_LIST_CCY")
	if listCcy == "" {
		listCcy = "BTC"
	}
	listRatio := os.Getenv("OKX_LIST_RATIO")
	if listRatio == "" {
		listRatio = "1"
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

	svc := c.NewTradingBotRecurringOrderAlgoService().
		StgyName(stgyName).
		RecurringList([]okx.TradingBotRecurringListItem{
			{Ccy: listCcy, Ratio: listRatio},
		}).
		Period(period).
		RecurringTime(recurringTime).
		TimeZone(timeZone).
		Amt(amt).
		InvestmentCcy(investmentCcy).
		TdMode(tdMode)

	if v := os.Getenv("OKX_RECURRING_DAY"); v != "" {
		svc.RecurringDay(v)
	}
	if v := os.Getenv("OKX_RECURRING_HOUR"); v != "" {
		svc.RecurringHour(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tradingBot recurring order placed: algoId=%s algoClOrdId=%s", ack.AlgoId, ack.AlgoClOrdId)
}
