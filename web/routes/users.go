package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/models"
)

func GetUserFromContext(c *gin.Context) *models.User {
	u, err := c.Get("user")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return nil
	}
	return u.(*models.User)
}
