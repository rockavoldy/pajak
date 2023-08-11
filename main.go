package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rockavoldy/pajak/kurs"

	"github.com/go-chi/chi/v5"
)

var HTTP_PORT = ":8080"

func main() {
	httpPortEnv := os.Getenv("HTTP_PORT")
	if httpPortEnv != "" {
		HTTP_PORT = httpPortEnv
	}

	r := chi.NewRouter()
	r.Mount("/kurs", kurs.Router())

	// s := gocron.NewScheduler(time.UTC)
	// s.Every(1).Week().Weekday(time.Tuesday).At("22:00:00").Do(func() {
	// 	log.Println("=== Start cron get Kurs Data ===")
	// 	log.Println("Cron run at: ", time.Now())

	// kursData := &KursData{}
	// getKursData(kursData)
	// err := kursData.CreateJson()
	// if err != nil {
	// 	log.Println(err)
	// }
	// })

	// s.StartImmediately()

	// s.StartAsync()
	log.Println("Listening on ", HTTP_PORT)
	if err := http.ListenAndServe(HTTP_PORT, r); err != nil {
		log.Panicln(err)
	}
}
