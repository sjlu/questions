package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"log"
	"net/http"
)

type Question struct {
	Question string
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("question")

	var questions []Question
	q.GetAll(c, &questions)

	log.Println(questions)

	json.NewEncoder(w).Encode(questions)
}

func AddQuestion(w http.ResponseWriter, r *http.Request) {
	vars := r.Form

	log.Println(vars)

	c := appengine.NewContext(r)
	q := Question{
		Question: r.FormValue("question"),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "question", nil), &q)
	if err != nil {
	}

	json.NewEncoder(w).Encode(q)
}
