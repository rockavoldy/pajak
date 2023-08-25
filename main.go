package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
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

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Week().Weekday(time.Tuesday).At("22:00:00").Do(func() {
		log.Println("=== Update kurs.json file ===")
		log.Println("Cron run at: ", time.Now())
		if err := kurs.UpdateKurs(); err != nil {
			log.Printf("err scheduler: %s", err)
		}
	})

	s.StartAsync()
	log.Println("Listening on ", HTTP_PORT)
	if err := http.ListenAndServe(HTTP_PORT, r); err != nil {
		log.Panicln(err)
	}
}
