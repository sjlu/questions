package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"io"
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

func GetQuestion(c appengine.Context, id int64) (*Question, error) {

	var question Question
	k := datastore.NewKey(c, "question", "", id, nil)
	err := datastore.Get(c, k, &question)
	if err != nil {
		return nil, err
	}

	question.Id = k.IntID()

	return &question, nil

}
