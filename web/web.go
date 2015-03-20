package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/questions", GetQuestions).Methods("GET")
	r.HandleFunc("/api/questions", AddQuestion).Methods("POST")
	http.Handle("/", r)
}
