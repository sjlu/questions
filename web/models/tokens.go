package models

import (
	"appengine"
	"appengine/datastore"
	"github.com/asaskevich/govalidator"
)

type Token struct {
	Id     int64 `json:"-" datastore:"-"`
	UserId int64 `json:"-" valid:"required"`
}

func (token *Token) key(c appengine.Context) *datastore.Key {
	if token.Id == 0 {
		return datastore.NewIncompleteKey(c, "Token", nil)
	}
	return datastore.NewKey(c, "Token", "", token.Id, nil)
}

func (token *Token) save(c appengine.Context) error {
	_, err := govalidator.ValidateStruct(token)
	if err != nil {
		return err
	}

	k, err := datastore.Put(c, token.key(c), token)
	if err != nil {
		return err
	}

	token.Id = k.IntID()
	return nil
}

func NewToken(c appengine.Context, user *User) (*Token, error) {
	var token Token
	token.UserId = user.Id

	err := token.save(c)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GetToken(c appengine.Context, id int64) (*Token, error) {
	var token Token
	token.Id = id

	k := token.key(c)
	err := datastore.Get(c, k, &token)
	if err != nil {
		return nil, err
	}

	token.Id = k.IntID()

	return &token, nil
}
