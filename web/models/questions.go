package models

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"io"
)

type Question struct {
	Id            int64      `json:"id" datastore:"-"`
	CategoryIds   []int64    `json:"category_ids"`
	Categories    []Category `json:"categories" datastore:"-"`
	UserId        int64      `json:"user_id"`
	User          *User      `json:"user" datastore:"-"`
	Question      string     `json:"question" valid:"required"`
	Answer1       string     `json:"answer_1" valid:"required"`
	Answer2       string     `json:"answer_2"`
	Answer3       string     `json:"answer_3"`
	Answer4       string     `json:"answer_4"`
	Answer5       string     `json:"answer_5"`
	Explanation   string     `json:"explanation" valid:"required"`
	CorrectAnswer string     `json:"correct_answer" valid:"required"`
	State         string     `json:"state" valid:"required"`
}

func (q *Question) key(c appengine.Context) *datastore.Key {
	if q.Id == 0 {
		return datastore.NewIncompleteKey(c, "Question", nil)
	}
	return datastore.NewKey(c, "Question", "", q.Id, nil)
}

func (q *Question) save(c appengine.Context) error {
	if q.State == "" {
		q.State = "new"
	}

	_, err := govalidator.ValidateStruct(q)
	if err != nil {
		return err
	}

	if q.State != "new" {
		return errors.New("state is invalid")
	}

	if q.CorrectAnswer != "answer_1" &&
		q.CorrectAnswer != "answer_2" &&
		q.CorrectAnswer != "answer_3" &&
		q.CorrectAnswer != "answer_4" &&
		q.CorrectAnswer != "answer_5" {
		return errors.New("correct_answer is invalid")
	}

	k, err := datastore.Put(c, q.key(c), q)
	if err != nil {
		return err
	}

	q.Id = k.IntID()
	return nil
}

func GetQuestions(c appengine.Context, userId int64) ([]Question, error) {
	q := datastore.NewQuery("Question")

	if userId != 0 {
		q = q.Filter("UserId =", userId)
	}

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

func NewQuestion(c appengine.Context, r io.ReadCloser, user *User) (*Question, error) {

	var question Question
	err := json.NewDecoder(r).Decode(&question)
	if err != nil {
		return nil, err
	}

	question.Id = 0
	question.UserId = user.Id

	err = question.save(c)
	if err != nil {
		return nil, err
	}

	return &question, nil

}

func GetQuestion(c appengine.Context, id int64) (*Question, error) {

	var question Question
	question.Id = id

	k := question.key(c)
	err := datastore.Get(c, k, &question)
	if err != nil {
		return nil, err
	}

	question.Id = k.IntID()

	if question.UserId != 0 {
		user, err := GetUser(c, question.UserId)
		if err != nil {
			return nil, err
		}
		question.User = user
	}

	if question.CategoryIds != nil {
		categories, err := GetCategoriesByIds(c, question.CategoryIds)
		if err != nil {
			return nil, err
		}
		question.Categories = categories
	}

	return &question, nil

}

func UpdateQuestion(c appengine.Context, id int64, r io.ReadCloser) (*Question, error) {

	var question Question
	question.Id = id

	k := question.key(c)
	err := datastore.Get(c, k, &question)
	if err != nil {
		return nil, err
	}

	var temp Question
	err = json.NewDecoder(r).Decode(&temp)
	if err != nil {
		return nil, err
	}

	question.Answer1 = temp.Answer1
	question.Answer2 = temp.Answer2
	question.Answer3 = temp.Answer3
	question.Answer4 = temp.Answer4
	question.Answer5 = temp.Answer5
	question.Question = temp.Question
	question.Explanation = temp.Explanation
	question.CategoryIds = temp.CategoryIds

	err = question.save(c)
	if err != nil {
		return nil, err
	}

	return &question, nil

}
