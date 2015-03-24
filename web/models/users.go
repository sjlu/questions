package models

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"io"
)

type User struct {
	Id    int64  `json:"id" datastore:"-"`
	Name  string `json:"name" valid:"required"`
	Email string `json:"email" valid:"required,email"`
	Role  string `json:"role"`
}

func (user *User) key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", "", user.Id, nil)
}

func (user *User) save(c appengine.Context) error {
	if user.Role == "" {
		user.Role = "user"
	}

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}

	k, err := datastore.Put(c, user.key(c), user)
	if err != nil {
		return err
	}

	user.Id = k.IntID()
	return nil
}

func NewUser(c appengine.Context, r io.ReadCloser) (*User, error) {
	var user User
	err := json.NewDecoder(r).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = user.save(c)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(c appengine.Context, id int64) (*User, error) {
	var user User
	user.Id = id

	k := user.key(c)
	err := datastore.Get(c, k, &user)
	if err != nil {
		return nil, err
	}

	user.Id = k.IntID()

	return &user, nil
}

func CreateOrUpdateUser(c appengine.Context, id int64, name string, email string) (*User, error) {
	user, err := GetUser(c, id)

	user.Email = email
	user.Name = name

	err = user.save(c)
	return user, err
}
