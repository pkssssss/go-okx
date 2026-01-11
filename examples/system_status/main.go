package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	state := os.Getenv("OKX_STATUS_STATE")

	c := okx.NewClient()
	svc := c.NewSystemStatusService()
	if state != "" {
		svc.State(state)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("count=%d\n", len(data))
	if len(data) > 0 {
		fmt.Printf("state=%s begin=%d end=%d title=%q\n", data[0].State, data[0].Begin, data[0].End, data[0].Title)
	}
}
