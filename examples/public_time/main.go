package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx"
)

func main() {
	c := okx.NewClient()

	st, err := c.NewPublicTimeService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ts=%d time=%s\n", st.TS, st.Time().Format("2006-01-02T15:04:05.000Z"))
}
