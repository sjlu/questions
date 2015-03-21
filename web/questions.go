package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	// "log"
	"github.com/gorilla/mux"
	"io"
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

func NewQuestion(c appengine.Context, r io.ReadCloser) (*Question, error) {

	var question Question
	err := json.NewDecoder(r).Decode(&question)
	if err != nil {
		return nil, err
	}

	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "question", nil), &question)
	if err != nil {
		return nil, err
	}

	question.Id = key.IntID()

	return &question, nil

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
