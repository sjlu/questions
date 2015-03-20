package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/questions", GetQuestions).Methods("GET")
	r.HandleFunc("/api/questions", AddQuestion).Methods("POST")
	r.HandleFunc("/api/questions/{id:[0-9]+}", GetQuestion).Methods("GET")
	http.Handle("/", r)
}
