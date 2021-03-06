package main

import (
	"appengine"
	"appengine/urlfetch"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	fb "github.com/huandu/facebook"
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

var SessionStore = sessions.NewCookieStore([]byte("yftcK6tjjW257QkwHuaqUHe8sj3s83Ky"))

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
	var config string = "env"
	if _, err := os.Stat("env.local"); os.IsNotExist(err) {
		config = "env"
	} else {
		if appengine.IsDevAppServer() {
			config = "env.local"
		} else {
			config = "env"
		}
	}

	err := godotenv.Load(config)
	if err != nil {
		log.Fatal("Problem loading configuration file.")
	}

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	//
	// middleware
	//
	r.Use(routes.GetUser)

	//
	// page routes
	//
	r.GET("/", func(c *gin.Context) {
		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/layout.tmpl", "templates/homepage.tmpl")))
		c.HTML(http.StatusOK, "layout", nil)
	})
	r.GET("/app", routes.RequiresUser, func(c *gin.Context) {
		session, err := SessionStore.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if session.Values["id"] == nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/layout.tmpl", "templates/app.tmpl")))
		c.HTML(http.StatusOK, "layout", nil)
	})

	//
	// authentication
	//
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
		session, err := SessionStore.Get(c.Request, "testable")
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
		session, err := SessionStore.Get(c.Request, "testable")
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
	login.POST("/token", func(c *gin.Context) {
		var json struct {
			FBAccessToken string `json:"FBAccessToken"`
		}
		c.Bind(&json)

		accessToken := json.FBAccessToken

		session := &fb.Session{}
		session.HttpClient = urlfetch.Client(appengine.NewContext(c.Request))
		session.SetAccessToken(accessToken)

		me, err := session.Get("/me", nil)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		id, err := strconv.ParseInt(me["id"].(string), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		name := me["name"].(string)
		email := me["email"].(string)

		user, err := models.CreateOrUpdateUser(appengine.NewContext(c.Request), id, name, email)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		token, err := models.NewToken(appengine.NewContext(c.Request), user)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": strconv.FormatInt(token.Id, 10)})
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

		session, err := SessionStore.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		session.Values["id"] = id
		session.Save(c.Request, c.Writer)

		c.Redirect(http.StatusTemporaryRedirect, "/app")
	})

	//
	// api endpoints
	//
	api := r.Group("/api")
	api.Use(routes.RequiresUser)

	routes.CategoryRouter(api.Group("/categories"))
	routes.QuestionRouter(api.Group("/questions"))

	//
	// http router init
	//
	http.Handle("/", r)
}
