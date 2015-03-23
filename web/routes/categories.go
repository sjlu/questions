package routes

import (
	"appengine"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		category, err := models.GetCategory(appengine.NewContext(c.Request), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, category)
	})

	router.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		category, err := models.RemoveCategory(appengine.NewContext(c.Request), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, category)
	})

	router.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		category, err := models.UpdateCategory(appengine.NewContext(c.Request), id, c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, category)
	})

}
