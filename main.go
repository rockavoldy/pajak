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
	"github.com/go-chi/httprate"
)

var HTTP_PORT = ":8080"

var (
	//go:embed templates
	res   embed.FS
	pages = map[string]string{
		"/":                "templates/index.html",
		"/dist/styles.css": "templates/dist/styles.css",
	}
)
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
	r.Use(httprate.Limit(
		10,
		10*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// view and assets
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	// })
	r.Get("/", indexHtml)
	r.Get("/dist/styles.css", stylesCss)

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
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}

	tmpl, err := template.ParseFS(res, page)
	if err != nil {
		log.Printf("page %s not found", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"userAgent": r.UserAgent(),
	}
	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}

func stylesCss(w http.ResponseWriter, r *http.Request) {
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}

	tmpl, err := template.ParseFS(res, page)
	if err != nil {
		log.Printf("page %s not found", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(http.StatusOK)

	data := map[string]interface{}{
		"userAgent": r.UserAgent(),
	}
	if err := tmpl.Execute(w, data); err != nil {
		return
	}
}
