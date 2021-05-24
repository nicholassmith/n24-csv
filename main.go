package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://www.24h-rennen.de/en/participants-2020/"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%s", doc.Children().Text())
	doc.Find("#content > div.container > div > div.col-md-12 > table").Children().Filter("#content > div.container > div > div.col-md-12 > table > tbody").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("node: %s", s.Contents().Text())
	})
}
