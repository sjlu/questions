package routes

import (
	"appengine"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"web/models"
)

func QuestionRouter(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		questions, err := models.GetQuestions(appengine.NewContext(c.Request))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, questions)
	})
	router.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		question, err := models.GetQuestion(appengine.NewContext(c.Request), id)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, question)
	})
	router.POST("/", func(c *gin.Context) {
		user := GetUserFromContext(c)

		question, err := models.NewQuestion(appengine.NewContext(c.Request), c.Request.Body, user)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, question)
	})
}
