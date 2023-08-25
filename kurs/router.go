package kurs

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", getKursHandler)
	r.Post("/update", updateKursHandler)

	return r
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeResponse(w http.ResponseWriter, status int, resp Response) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, status int, err error) {
	resp := Response{
		Message: err.Error(),
		Data:    nil,
	}
	writeResponse(w, status, resp)
}

func getKursHandler(w http.ResponseWriter, r *http.Request) {
	kursData, err := loadKurs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	status := http.StatusOK
	resp := Response{
		Message: http.StatusText(status),
		Data:    kursData,
	}
	writeResponse(w, status, resp)
}

func updateKursHandler(w http.ResponseWriter, r *http.Request) {
	err := UpdateKurs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	kursData, err := loadKurs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	status := http.StatusOK
	resp := Response{
		Message: http.StatusText(status),
		Data:    kursData,
	}
	writeResponse(w, status, resp)
}
