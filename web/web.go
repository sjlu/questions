package web

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/objx"
	"io"
	"net/http"
	"strconv"
)

func init() {

	gomniauth.SetSecurityKey("hlhi23o9fhlASdfaSDf078923oifsASDFAsdf8973r28y2y8")
	gomniauth.WithProviders(
		facebook.New(
			"1616407611915270",
			"6c0555b73ee505405926d83ff4f8ba7c",
			"http://localhost:8080/login/callback",
		),
	)

	r := mux.NewRouter()
	r.HandleFunc("/api/questions", getQuestions).Methods("GET")
	r.HandleFunc("/api/questions", newQuestion).Methods("POST")
	r.HandleFunc("/api/questions/{id:[0-9]+}", getQuestion).Methods("GET")
	r.HandleFunc("/api/categories", getCategories).Methods("GET")
	r.HandleFunc("/api/categories", newCategory).Methods("POST")
	r.HandleFunc("/login", loginHandler())
	r.HandleFunc("/login/callback", callbackHandler())
	http.Handle("/", r)

}

func loginHandler() http.HandlerFunc {
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		t := new(urlfetch.Transport)
		t.Context = c
		common.SetRoundTripper(t)

		state := gomniauth.NewState("after", "success")
		authUrl, err := provider.GetBeginAuthURL(state, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func callbackHandler() http.HandlerFunc {
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		t := new(urlfetch.Transport)
		t.Context = c
		common.SetRoundTripper(t)

		omap, err := objx.FromURLQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		creds, err := provider.CompleteAuth(omap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := fmt.Sprintf("%#v", user)
		io.WriteString(w, data)

	}
}

func getQuestions(w http.ResponseWriter, r *http.Request) {

	questions, err := GetQuestions(appengine.NewContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(questions)

}

func newQuestion(w http.ResponseWriter, r *http.Request) {

	question, err := NewQuestion(appengine.NewContext(r), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(question)

}

func getQuestion(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	question, err := GetQuestion(appengine.NewContext(r), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(question)

}

func newCategory(w http.ResponseWriter, r *http.Request) {

	category, err := NewCategory(appengine.NewContext(r), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(category)

}

func getCategories(w http.ResponseWriter, r *http.Request) {

	topics, err := GetCategories(appengine.NewContext(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(topics)

}
