package routes

import (
	"appengine"
	"github.com/gin-gonic/gin"
	"net/http"
	"web/models"
)

func CategoryRouter(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		categories, err := models.GetCategories(appengine.NewContext(c.Request))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, categories)
	})
	router.POST("/", func(c *gin.Context) {
		category, err := models.NewCategory(appengine.NewContext(c.Request), c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, category)
	})
}
