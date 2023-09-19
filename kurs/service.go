package kurs

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rockavoldy/pajak/currency"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	ErrMissingValidDate = errors.New("validation date is missing")
)

func loadKurs() (KursData, error) {
	// no need to check the timestamp, since it will be updated regularly with cron
	currDir, err := os.Getwd()
	if err != nil {
		return KursData{}, err
	}

	content, err := os.ReadFile(currDir + "/dist/kurs.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			UpdateKurs()
		}
	}

	var kursData KursData
	err = json.Unmarshal(content, &kursData)
	if err != nil {
		return KursData{}, err
	}
	return kursData, nil
}

func UpdateKurs() error {
	kursData, err := getNewKurs()
	if err != nil {
		return err
	}

	err = createJsonKurs(kursData)
	if err != nil {
		return err
	}

	return nil
}

func getNewKurs() (KursData, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("fiskal.kemenkeu.go.id", "kemenkeu.go.id"),
		colly.UserAgent("Googlebot/2.1 (http://www.googlebot.com/bot.html)"),
	)

	kursData := NewKursData()

	c.OnHTML("p.text-muted", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.Text)
		text = strings.TrimPrefix(text, "Tanggal Berlaku: ")
		textSplitted := strings.Split(text, " - ")
		kursData.ValidFrom = textSplitted[0]
		kursData.ValidTo = textSplitted[1]
	})

	c.OnHTML("table tbody", func(h *colly.HTMLElement) {
		h.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
			var name string
			var symbol string
			var value float64
			var changes float64
			s.Find("td").Each(func(index int, el *goquery.Selection) {
				if index == 1 {
					name = el.Find("span.hidden-xs").Text()
					symbol = el.Find("span.visible-xs-inline").Text()
				} else if index == 2 {
					val := el.Find("div.ml-5").Text()
					val = strings.ReplaceAll(val, ".", "")
					val = strings.ReplaceAll(val, ",", "")

					valu, _ := strconv.ParseInt(val, 10, 64)
					if strings.Compare(symbol, "JPY") == 0 {
						value = float64(valu) / 100
					} else {
						value = float64(valu)
					}
				} else if index == 3 {
					val := strings.TrimSpace(el.Text())
					val = strings.ReplaceAll(val, ".", "")
					val = strings.ReplaceAll(val, ",", "")

					change, _ := strconv.ParseInt(val, 10, 32)
					if strings.Compare(symbol, "JPY") == 0 {
						changes = float64(change) / 100
					} else {
						changes = float64(change)
					}
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

func createJsonKurs(kursData KursData) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.MkdirAll(currDir+"/dist", 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	kursDataFile, err := os.OpenFile(currDir+"/dist/kurs.json", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	defer kursDataFile.Close()

	kursDataJson, err := json.Marshal(kursData)
	if err != nil {
		return err
	}

	err = os.WriteFile(currDir+"/dist/kurs.json", kursDataJson, 0644)
	if err != nil {
		return err
	}

	return nil
}
