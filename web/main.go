package main

import (
	"appengine"
	"appengine/urlfetch"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"web/models"
	"web/routes"
)

var store = sessions.NewCookieStore([]byte("yftcK6tjjW257QkwHuaqUHe8sj3s83Ky"))

func loadConfig(file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return
	}
	err := godotenv.Load(file)
	if err != nil {
		log.Fatal("Problem loading configuration file.")
	}
}

func init() {
	var config string = ".env"
	if _, err := os.Stat(".env.local"); os.IsNotExist(err) {
		config = ".env"
	} else {
		if appengine.IsDevAppServer() {
			config = ".env.local"
		} else {
			config = ".env"
		}
	}

	err := godotenv.Load(config)
	if err != nil {
		log.Fatal("Problem loading configuration file.")
	}

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/app", func(c *gin.Context) {
		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["id"] == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/layout.tmpl", "templates/app.tmpl")))
		c.HTML(200, "layout", nil)
	})

	gomniauth.SetSecurityKey("mJ8zwRBQZqvakN2BT6CuVKQD8gxYXW8X")
	gomniauth.WithProviders(
		facebook.New(
			os.Getenv("FACEBOOK_APP_ID"),
			os.Getenv("FACEBOOK_APP_SECRET"),
			os.Getenv("FACEBOOK_CALLBACK_URL"),
		),
	)
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}

	r.GET("/logout", func(c *gin.Context) {
		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		session.Values["id"] = nil
		session.Save(c.Request, c.Writer)

		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	login := r.Group("/login")
	login.GET("/", func(c *gin.Context) {
		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["id"] != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/app")
			return
		}

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

		id, err := strconv.ParseInt(user.IDForProvider("facebook"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		name := user.Name()
		email := user.Email()

		_, err = models.CreateOrUpdateUser(appengine.NewContext(c.Request), id, name, email)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		session, err := store.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		session.Values["id"] = id
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

		if session.Values["id"] == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id := session.Values["id"].(int64)
		user, err := models.GetUser(appengine.NewContext(c.Request), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Set("user", user)
	})

	routes.CategoryRouter(api.Group("/categories"))
	routes.QuestionRouter(api.Group("/questions"))

	http.Handle("/", r)
}
