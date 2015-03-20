package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	// "log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Question struct {
	Id       int64  `json:"id" datastore:"-"`
	Question string `json:"question"`
}

func GetQuestions(c appengine.Context) ([]Question, error) {
	q := datastore.NewQuery("question")

	var questions []Question
	keys, err := q.GetAll(c, &questions)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(questions); i++ {
		questions[i].Id = keys[i].IntID()
	}

	return questions, nil
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

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k := datastore.NewKey(c, "question", "", id, nil)
	var question Question
	err = datastore.Get(c, k, &question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	question.Id = id

	json.NewEncoder(w).Encode(question)
}
