package main

import (
	"context"
	"log"
	"os"
	"time"

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

	timeout := 30 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
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

	msgCh := make(chan okx.EconomicCalendarEvent, 1)
	errCh := make(chan error, 1)

	ws := c.NewWSBusinessPrivate(okx.WithWSTypedHandlerAsync(1024))
	ws.OnEconomicCalendar(func(event okx.EconomicCalendarEvent) {
		select {
		case msgCh <- event:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}); err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	subCtx, cancelSub := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelEconomicCalendar}); err != nil {
		log.Fatal(err)
	}

	select {
	case e := <-msgCh:
		log.Printf("economic-calendar: calendarId=%s region=%s category=%s event=%q importance=%s date=%d refDate=%d actual=%s forecast=%s previous=%s ts=%d",
			e.CalendarId, e.Region, e.Category, e.Event, e.Importance, e.Date, e.RefDate, e.Actual, e.Forecast, e.Previous, e.TS)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no economic-calendar push within %s", timeout)
	}

	cancel()
	<-ws.Done()
}
