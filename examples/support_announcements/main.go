package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	annType := os.Getenv("OKX_ANN_TYPE")
	page := os.Getenv("OKX_PAGE")
	lang := os.Getenv("OKX_ACCEPT_LANGUAGE")
	if lang == "" {
		lang = "zh-CN"
	}

	c := okx.NewClient()
	data, err := c.NewSupportAnnouncementsService().
		AnnType(annType).
		Page(page).
		AcceptLanguage(lang).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("totalPage=%s count=%d\n", data.TotalPage, len(data.Details))
	if len(data.Details) > 0 {
		v := data.Details[0]
		fmt.Printf("annType=%s title=%q url=%s pTime=%d businessPTime=%d\n", v.AnnType, v.Title, v.URL, v.PTime, v.BusinessPTime)
	}
}
