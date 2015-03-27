package routes

import (
	"appengine"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
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
	session, err := SessionStore.Get(c.Request, "testable")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	id := session.Values["id"].(int64)
	user, err := models.GetUser(appengine.NewContext(c.Request), id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
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
