package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()
	data, err := c.NewSupportAnnouncementTypesService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count=%d\n", len(data))
	for i, v := range data {
		if i >= 10 {
			break
		}
		fmt.Printf("%s: %s\n", v.AnnType, v.AnnTypeDesc)
	}
}
