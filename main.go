package main

import (
	"appengine"
	"appengine/urlfetch"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/objx"
	"html/template"
	"models"
	"net/http"
	"routes"
	"strconv"
)

var store = sessions.NewCookieStore([]byte("yftcK6tjjW257QkwHuaqUHe8sj3s83Ky"))

func init() {

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/app", func(c *gin.Context) {
		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["email"] == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/layout.tmpl", "templates/app.tmpl")))
		c.HTML(200, "layout", nil)
	})

	gomniauth.SetSecurityKey("mJ8zwRBQZqvakN2BT6CuVKQD8gxYXW8X")
	gomniauth.WithProviders(
		facebook.New(
			"1616407611915270",
			"6c0555b73ee505405926d83ff4f8ba7c",
			"http://localhost:8080/login/callback",
		),
	)
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}

	login := r.Group("/login")
	login.GET("/", func(c *gin.Context) {
		t := new(urlfetch.Transport)
		t.Context = appengine.NewContext(c.Request)
		common.SetRoundTripper(t)

		state := gomniauth.NewState("after", "success")
		authUrl, err := provider.GetBeginAuthURL(state, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, authUrl)
	})
	login.GET("/callback", func(c *gin.Context) {
		t := new(urlfetch.Transport)
		t.Context = appengine.NewContext(c.Request)
		common.SetRoundTripper(t)

		omap, err := objx.FromURLQuery(c.Request.URL.RawQuery)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		creds, err := provider.CompleteAuth(omap)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		session.Values["email"] = user.Email()
		session.Save(c.Request, c.Writer)

		c.JSON(http.StatusOK, user.Data())

		id, err := strconv.ParseInt(user.IDForProvider("facebook"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		name := user.Name()
		email := user.Email()

		_, err = models.CreateUser(appengine.NewContext(c.Request), id, name, email)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, "/app")
	})

	api := r.Group("/api")
	api.Use(func(c *gin.Context) {
		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["email"] == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	})

	routes.CategoryRouter(api.Group("/categories"))
	routes.QuestionRouter(api.Group("/questions"))

	http.Handle("/", r)
}
