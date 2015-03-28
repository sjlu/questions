package routes

import (
	"appengine"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
	"web/models"
)

var SessionStore = sessions.NewCookieStore([]byte("yftcK6tjjW257QkwHuaqUHe8sj3s83Ky"))

func GetUserFromContext(c *gin.Context) *models.User {
	u, err := c.Get("user")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil
	}
	return u.(*models.User)
}

func GetUser(c *gin.Context) {
	var id int64

	authToken := c.Request.Header.Get("authentication-token")
	if authToken != "" {
		tokenId, err := strconv.ParseInt(authToken, 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		token, err := models.GetToken(appengine.NewContext(c.Request), tokenId)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		id = token.UserId
	} else {
		session, err := SessionStore.Get(c.Request, "testable")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		idString := session.Values["id"]
		if idString != nil {
			var ok bool
			id, ok = session.Values["id"].(int64)
			if !ok {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	var user *models.User
	if id != 0 {
		var err error
		user, err = models.GetUser(appengine.NewContext(c.Request), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	if user != nil {
		c.Set("user", user)
	}
}

func RequiresUser(c *gin.Context) {
	user := GetUserFromContext(c)
	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RequiresAdmin(c *gin.Context) {
	user := GetUserFromContext(c)
	if user.Role != "admin" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
