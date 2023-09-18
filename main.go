package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rockavoldy/pajak/kurs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var HTTP_PORT = ":8080"

//go:embed templates/*
var res embed.FS
var t *template.Template

func main() {
	t = template.Must(template.ParseFS(res, "templates/index.html"))
	template.Must(t.ParseFS(res, "templates/dist/*"))

	httpPortEnv := os.Getenv("HTTP_PORT")
	if httpPortEnv != "" {
		HTTP_PORT = httpPortEnv
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// view and assets
	r.Get("/", indexHtml)
	r.Get("/styles.css", stylesCss)

	// kurs API serve
	r.Mount("/kurs", kurs.Router())

	// Scheduler to run every Wednesday 5AM UTC+7 (new kurs updated every Wednesday)
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

func indexHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	t.ExecuteTemplate(w, "index.html", nil)
}

func stylesCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/css")
	t.ExecuteTemplate(w, "styles.css", nil)
}
