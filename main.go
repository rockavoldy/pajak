package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
	_ "github.com/mattn/go-sqlite3"
)

type Currency struct {
	Currency string `json:"currency"`
	Symbol   string `json:"symbol"`
	Value    int    `json:"value"`   // This value need to be divided with per 100
	Changes  int    `json:"changes"` // This changes should be divided with per 100
}

type KursData struct {
	UpdatedAt  uint       `json:"updated_at"`
	ValidFrom  string     `json:"valid_from"`
	ValidTo    string     `json:"valid_to"`
	Currencies []Currency `json:"currencies"`
}

func (c *Currency) FillCurrency(index int, el *goquery.Selection) {
	if index == 1 {
		c.Currency = el.Find("span.hidden-xs").Text()
		c.Symbol = el.Find("span.visible-xs-inline").Text()
	} else if index == 2 {
		val := el.Find("div.ml-5").Text()
		val = strings.ReplaceAll(val, ".", "")
		val = strings.ReplaceAll(val, ",", "")

		value, _ := strconv.ParseInt(val, 10, 64)
		c.Value = int(value)
	} else if index == 3 {
		val := strings.TrimSpace(el.Text())
		val = strings.ReplaceAll(val, ".", "")
		val = strings.ReplaceAll(val, ",", "")

		changes, _ := strconv.ParseInt(val, 10, 32)
		c.Changes = int(changes)
	}
}

func (k *KursData) FillKursData(index int, el *goquery.Selection) {
	currency := &Currency{}
	el.Find("td").Each(currency.FillCurrency)
	k.Currencies = append(k.Currencies, *currency)
}

func (k *KursData) getChecksum() string {
	kurs := &KursData{
		Currencies: k.Currencies,
		ValidFrom:  k.ValidFrom,
		ValidTo:    k.ValidTo,
	}
	jsonCurrencies, err := json.Marshal(kurs)
	if err != nil {
		log.Panic(err)
		return ""
	}

	md5Checksum := md5.Sum(jsonCurrencies)

	return hex.EncodeToString(md5Checksum[:])
}

func getKursData(kursData *KursData) error {
	var err error

	c := colly.NewCollector(
		colly.AllowedDomains("fiskal.kemenkeu.go.id", "kemenkeu.go.id"),
		colly.UserAgent("Googlebot/2.1 (http://www.googlebot.com/bot.html)"),
	)

	c.OnHTML("table tbody", func(h *colly.HTMLElement) {
		kursData.UpdatedAt = uint(time.Now().Unix())
		h.DOM.Find("tr").Each(kursData.FillKursData)
	})

	c.OnHTML("p.text-muted", func(h *colly.HTMLElement) {
		text := strings.TrimSpace(h.Text)
		text = strings.TrimPrefix(text, "Tanggal Berlaku: ")
		textSplitted := strings.Split(text, " - ")
		kursData.ValidFrom = textSplitted[0]
		kursData.ValidTo = textSplitted[1]

	})

	c.OnError(func(r *colly.Response, e error) {
		err = e
	})

	c.Visit("https://fiskal.kemenkeu.go.id/informasi-publik/kurs-pajak")

	if err != nil {
		return err
	}

	return nil
}

func (k *KursData) Insert(db *sql.DB) error {
	sql, args, err := sq.Insert("kurs").
		Columns("valid_from", "valid_to", "updated_at").
		Values(k.ValidFrom, k.ValidTo, k.UpdatedAt).
		ToSql()
	if err != nil {
		return err
	}

	tx, _ := db.Begin()
	kurs, err := tx.Exec(sql, args...)
	if err != nil {
		return err
	}

	kurs_id, err := kurs.LastInsertId()
	if err != nil {
		return err
	}

	for _, currency := range k.Currencies {
		sql, args, _ = sq.Insert("currency").
			Columns("kurs_id", "currency", "symbol", "value", "changes").
			Values(kurs_id, currency.Currency, currency.Symbol, currency.Value, currency.Changes).ToSql()
		tx.Exec(sql, args...)
	}

	newChecksum := k.getChecksum()
	sql, args, _ = sq.Insert("kurs_checksum").Columns("kurs_id", "checksum").Values(kurs_id, newChecksum).ToSql()
	tx.Exec(sql, args...)
	tx.Commit()

	return nil
}

func (k *KursData) checkDuplicate(db *sql.DB) bool {
	newChecksum := k.getChecksum()
	if newChecksum == "" {
		return false
	}
	sql, args, _ := sq.Select("*").From("kurs_checksum").Where(sq.Eq{"checksum": newChecksum}).Limit(1).ToSql()

	query, err := db.Query(sql, args...)
	if err != nil {
		log.Println(err)
	}

	if query.Next() {
		return true
	}

	return false
}

func (k *KursData) CreateJson() error {
	dirs, err := os.ReadDir("dist")

	var currencies []map[string]interface{}
	for _, v := range k.Currencies {
		value := fmt.Sprintf("%.2f", float64(v.Value)/100)
		changes := fmt.Sprintf("%.2f", float64(v.Changes)/100)
		if strings.Compare(v.Symbol, "JPY") == 0 {
			value = fmt.Sprintf("%.2f", float64(v.Value)/100/100)
			changes = fmt.Sprintf("%.2f", float64(v.Changes)/100/100)
		}
		currency := &map[string]interface{}{
			"currency": v.Currency,
			"symbol":   v.Symbol,
			"value":    value,
			"changes":  changes,
		}
		currencies = append(currencies, *currency)
	}
	kursPrint := &map[string]interface{}{
		"valid_from": k.ValidFrom,
		"valid_to":   k.ValidTo,
		"updated_at": k.UpdatedAt,
		"currencies": currencies,
	}

	if err != nil {
		err = os.Mkdir("dist", 0755)
		if err != nil {
			return err
		}
		dirs, _ = os.ReadDir("dist")
	}

	os.Chdir("dist")

	if len(dirs) == 0 {
		os.Create("kurs.json")
	}

	if len(dirs) > 0 && dirs[0].Name() != "kurs.json" {
		os.Create("kurs.json")
	}

	jsonData, _ := json.Marshal(kursPrint)
	err = os.WriteFile("kurs.json", jsonData, 0666)
	os.Chdir("..")
	if err != nil {
		return err
	}

	return nil
}

// ResponseJSON. Will response to http with json data
func ResponseJSON(rw http.ResponseWriter, data interface{}, status int) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if status == 0 {
		status = http.StatusOK
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	rw.Write([]byte(dataJson))
	return nil
}

// PrintResponse, to give response with single data (like get, create, update, or delete)
type PrintResponse struct {
	ResponseCode int                    `json:"response_code"`
	Message      string                 `json:"message"`
	Data         map[string]interface{} `json:"data"`
}

func NewRouter() *chi.Mux {
	router := chi.NewMux()
	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		data := PrintResponse{
			ResponseCode: http.StatusOK,
			Message:      "Kurs Pajak API: Diterbitkan oleh Kemenkeu setiap hari rabu.",
		}

		ResponseJSON(rw, data, data.ResponseCode)
	})
	router.MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
		data := PrintResponse{
			ResponseCode: http.StatusMethodNotAllowed,
			Message:      "405 Method not allowed",
		}

		ResponseJSON(rw, data, data.ResponseCode)
	})
	router.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		data := PrintResponse{
			ResponseCode: http.StatusNotFound,
			Message:      "404 Not found",
		}

		ResponseJSON(rw, data, data.ResponseCode)
	})

	return router
}

var HTTP_PORT = ":8080"

func main() {
	httpPortEnv := os.Getenv("HTTP_PORT")
	if httpPortEnv != "" {
		HTTP_PORT = httpPortEnv
	}
	router := NewRouter()

	router.Get("/kurs", func(rw http.ResponseWriter, r *http.Request) {
		currDir, err := os.Getwd()
		data := PrintResponse{}
		if err != nil {
			data = PrintResponse{
				ResponseCode: http.StatusInternalServerError,
				Message:      err.Error(),
			}
		} else {
			content, _ := os.ReadFile(currDir + "/dist/kurs.json")
			var payload map[string]interface{}
			json.Unmarshal(content, &payload)
			data = PrintResponse{
				ResponseCode: http.StatusOK,
				Message:      "Success",
				Data:         payload,
			}
		}

		ResponseJSON(rw, data, data.ResponseCode)
	})

	router.Post("/update-kurs", func(rw http.ResponseWriter, r *http.Request) {
		kursData := &KursData{}
		getKursData(kursData)
		data := PrintResponse{
			ResponseCode: http.StatusOK,
			Message:      "Success update data",
		}

		err := kursData.CreateJson()
		if err != nil {
			log.Println(err)
			data = PrintResponse{
				ResponseCode: http.StatusBadRequest,
				Message:      "Failed update data",
			}
		}

		ResponseJSON(rw, data, data.ResponseCode)
	})

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Week().Weekday(time.Tuesday).At("22:00:00").Do(func() {
		log.Println("=== Start cron get Kurs Data ===")
		log.Println("Cron run at: ", time.Now())

		kursData := &KursData{}
		getKursData(kursData)
		err := kursData.CreateJson()
		if err != nil {
			log.Println(err)
		}
	})

	s.StartImmediately().StartAsync()

	s.StartAsync()
	log.Println("Listening on ", HTTP_PORT)
	if err := http.ListenAndServe(HTTP_PORT, router); err != nil {
		log.Panicln(err)
	}
}
