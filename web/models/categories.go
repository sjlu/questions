package models

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"io"
)

type Category struct {
	Id   int64  `json:"id" datastore:"-"`
	Name string `json:"name" valid:"required"`
}

func (category *Category) key(c appengine.Context) *datastore.Key {
	if category.Id == 0 {
		return datastore.NewIncompleteKey(c, "Category", nil)
	}
	return datastore.NewKey(c, "Category", "", category.Id, nil)
}

func (category *Category) save(c appengine.Context) error {
	_, err := govalidator.ValidateStruct(category)
	if err != nil {
		return err
	}

	k, err := datastore.Put(c, category.key(c), category)
	if err != nil {
		return err
	}

	category.Id = k.IntID()
	return nil
}

func GetCategories(c appengine.Context) ([]Category, error) {
	q := datastore.NewQuery("Category").Order("Name")

	var categories []Category
	keys, err := q.GetAll(c, &categories)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(categories); i++ {
		categories[i].Id = keys[i].IntID()
	}

	return categories, nil
}

func GetCategoriesByIds(c appengine.Context, ids []int64) ([]Category, error) {
	var keys []*datastore.Key

	for _, id := range ids {
		keys = append(keys, datastore.NewKey(c, "Category", "", id, nil))
	}

	categories := make([]Category, len(keys))
	err := datastore.GetMulti(c, keys, categories)
	for i := 0; i < len(categories); i++ {
		categories[i].Id = keys[i].IntID()
	}
	if err != nil {
		if multiErr, ok := err.(appengine.MultiError); ok {
			for i, e := range multiErr {
				if e == datastore.ErrNoSuchEntity {
					categories = append(categories[:i], categories[i+1:]...)
				} else if e != nil {
					return nil, multiErr
				}
			}
		}
	}

	return categories, nil
}

func GetCategory(c appengine.Context, id int64) (*Category, error) {
	var category Category
	category.Id = id

	k := category.key(c)
	err := datastore.Get(c, k, &category)
	if err != nil {
		return nil, err
	}

	category.Id = k.IntID()

	return &category, nil
}

func NewCategory(c appengine.Context, r io.ReadCloser) (*Category, error) {

	var category Category
	err := json.NewDecoder(r).Decode(&category)
	if err != nil {
		return nil, err
	}

	category.Id = 0

	err = category.save(c)
	if err != nil {
		return nil, err
	}

	return &category, nil

}

func RemoveCategory(c appengine.Context, id int64) (*Category, error) {

	category, err := GetCategory(c, id)
	if err != nil {
		return nil, err
	}

	err = datastore.Delete(c, category.key(c))
	if err != nil {
		return nil, err
	}

	return category, nil

}

func UpdateCategory(c appengine.Context, id int64, r io.ReadCloser) (*Category, error) {

	var category Category
	category.Id = id

	k := category.key(c)
	err := datastore.Get(c, k, &category)
	if err != nil {
		return nil, err
	}

	var cat Category
	err = json.NewDecoder(r).Decode(&cat)
	if err != nil {
		return nil, err
	}

	category.Name = cat.Name

	err = category.save(c)
	if err != nil {
		return nil, err
	}

	return &category, nil

}
