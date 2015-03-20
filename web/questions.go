package web

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
)

type Question struct {
	Question string
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	q := datastore.NewQuery("Question")

	var questions []Question
	q.GetAll(c, &questions)

	json.NewEncoder(w).Encode(questions)
}

func AddQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "got here")
}
