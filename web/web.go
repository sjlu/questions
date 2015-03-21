package web

import (
	"appengine"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/api/questions", getQuestions).Methods("GET")
	r.HandleFunc("/api/questions", newQuestion).Methods("POST")
	r.HandleFunc("/api/questions/{id:[0-9]+}", getQuestion).Methods("GET")
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

func getQuestion(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	question, err := GetQuestion(appengine.NewContext(r), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(question)

}
