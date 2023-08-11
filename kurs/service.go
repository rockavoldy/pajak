package kurs

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/rockavoldy/pajak/currency"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrMissingValidDate = errors.New("validation date is missing")
)

func getKurs(ctx context.Context) (KursData, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("fiskal.kemenkeu.go.id", "kemenkeu.go.id"),
		colly.UserAgent("Googlebot/2.1 (http://www.googlebot.com/bot.html)"),
	)

	kursData := NewKursData()

	c.OnHTML("p.text-muted", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.Text)
		text = strings.TrimPrefix(text, "Tanggal Berlaku: ")
		textSplitted := strings.Split(text, " - ")
		kursData.ValidFrom = null.StringFrom(textSplitted[0])
		kursData.ValidTo = null.StringFrom(textSplitted[1])
	})

	c.OnHTML("table tbody", func(h *colly.HTMLElement) {
		h.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
			var name string
			var symbol string
			var value int
			var changes int
			s.Find("td").Each(func(index int, el *goquery.Selection) {
				if index == 1 {
					name = el.Find("span.hidden-xs").Text()
					symbol = el.Find("span.visible-xs-inline").Text()
				} else if index == 2 {
					val := el.Find("div.ml-5").Text()
					val = strings.ReplaceAll(val, ".", "")
					val = strings.ReplaceAll(val, ",", "")

					valu, _ := strconv.ParseInt(val, 10, 64)
					value = int(valu)
				} else if index == 3 {
					val := strings.TrimSpace(el.Text())
					val = strings.ReplaceAll(val, ".", "")
					val = strings.ReplaceAll(val, ",", "")

					change, _ := strconv.ParseInt(val, 10, 32)
					changes = int(change)
				}
			})

			currency, err := currency.NewCurrency(name, symbol, value, changes)
			if err != nil {
				log.Println(err)
			}
			kursData.Currencies = append(kursData.Currencies, currency)
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("(%s): %s\n", r.Request.URL, err)
	})

	c.Visit("https://fiskal.kemenkeu.go.id/informasi-publik/kurs-pajak")

	return kursData, nil
}
