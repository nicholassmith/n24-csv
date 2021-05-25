package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type EntryNode struct {
	EntryNumber  int    `json:"entryNumber"`
	Class        string `json:"class"`
	Team         string `json:"team"`
	Manufacturer string `json:"manufacturer"`
	Car          string `json:"car"`
}

const (
	entryNumberRow     = 0
	entryClassRow      = 2
	entryTeamRow       = 3
	entryCarDetailsRow = 4
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

	entryList := make(map[string][]EntryNode)

	doc.Find("#content > div.container > div > div.col-md-12 > table").Children().Filter("#content > div.container > div > div.col-md-12 > table > tbody").Each(func(i int, s *goquery.Selection) {
		s.Children().Each(func(x int, node *goquery.Selection) {
			entryNode := new(EntryNode)
			node.Children().Each(func(y int, elem *goquery.Selection) {
				switch {
				case y == entryNumberRow:
					entry, _ := strconv.Atoi(elem.Text())
					entryNode.EntryNumber = entry
				case y == entryClassRow:
					entryNode.Class = elem.Text()
				case y == entryTeamRow:
					team := elem.Find("b").Text()
					entryNode.Team = team
				case y == entryCarDetailsRow:
					line, _ := elem.Html()
					split := strings.Split(line, "<br/>")
					entryNode.Manufacturer = split[0]
					entryNode.Car = split[1]
				}
			})
			entryList[entryNode.Class] = append(entryList[entryNode.Class], *entryNode)
		})
	})

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.Encode(entryList)

	fmt.Printf("%s", &buf)
}
