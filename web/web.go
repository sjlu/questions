package web

import (
	"appengine"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/questions", getQuestions).Methods("GET")
	r.HandleFunc("/api/questions", newQuestion).Methods("POST")
	r.HandleFunc("/api/questions/{id:[0-9]+}", GetQuestion).Methods("GET")
	http.Handle("/", r)
}

func getQuestions(w http.ResponseWriter, r *http.Request) {

	questions, err := GetQuestions(appengine.NewContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(questions)

}

func newQuestion(w http.ResponseWriter, r *http.Request) {

	question, err := NewQuestion(appengine.NewContext(r), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(question)

}
