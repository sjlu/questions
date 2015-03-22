package web

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
	"net/http"
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
	api.GET("/questions", func(c *gin.Context) {
		questions, err := GetQuestions(appengine.NewContext(c.Request))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, questions)
	})
	api.GET("/questions/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		question, err := GetQuestion(appengine.NewContext(c.Request), id)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, question)
	})
	api.POST("/questions", func(c *gin.Context) {
		question, err := NewQuestion(appengine.NewContext(c.Request), c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, question)
	})

	api.GET("/categories", func(c *gin.Context) {
		categories, err := GetCategories(appengine.NewContext(c.Request))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, categories)
	})
	api.POST("/categories", func(c *gin.Context) {
		category, err := NewCategory(appengine.NewContext(c.Request), c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, category)
	})

	http.Handle("/", r)
}
