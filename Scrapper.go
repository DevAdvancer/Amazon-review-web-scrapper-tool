package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strings"
)

func main() {
	// Product id = B0D9LJH5F5 for example
	var productId string
	fmt.Print("Enter the Id of Product from Amazon to Scrape: ")
	fmt.Scan(&productId)
	getProduct(productId)
}

func getProduct(asin string) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered("https://www.amazon.in/gp/product/"+asin+"?ref=ox_sc_act_title_4&smid=A1V7ZM32AEQ8C&th=1", g.Opt.ParseFunc)
		},
		ParseFunc: reviewParse,
		Exporters: []export.Exporter{&export.JSONLine{FileName: "text.json"}},
	}).Start()
}

func reviewParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div[data-hook=review]").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a[data-hook=review-title]").Text()
		cleanTitle := strings.TrimSpace(title)

		g.Exports <- map[string]interface{}{
			"title": cleanTitle,
			"date":  s.Find("span[data-hook=review-date]").Text(),
			"body":  s.Find("span[data-hook=review-body] span").Text(),
		}
	})

	if href, ok := r.HTMLDoc.Find("li.a-last a").Attr("href"); ok {
		g.GetRendered(r.JoinURL(href), reviewParse)
	}
}
