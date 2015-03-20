package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	// "log"
	"net/http"
)

type Question struct {
	Id       int64  `json:"id" datastore:"-"`
	Question string `json:"question"`
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("question")

	var questions []Question
	keys, err := q.GetAll(c, &questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(questions); i++ {
		questions[i].Id = keys[i].IntID()
	}

	json.NewEncoder(w).Encode(questions)
}

func AddQuestion(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var q Question
	err := json.NewDecoder(r.Body).Decode(&q)

	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "question", nil), &q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(q)
}