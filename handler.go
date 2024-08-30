package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	// add a new schedule
	mux.Post("/schedule/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	// get all schedules
	mux.Get("/schedule/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	// delete a schedule
	mux.Delete("/schedule/{userId}/{scheduleId}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	// add a user
	mux.Post("/user", func(w http.ResponseWriter, r *http.Request) {})
	// delete a user
	mux.Delete("/user", func(w http.ResponseWriter, r *http.Request) {})
	return mux
}

func AddScheduleHandler(w http.ResponseWriter, r *http.Request) {

}

func ListScheduleHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {

}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

}
