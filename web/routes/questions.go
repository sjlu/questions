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
		user := GetUserFromContext(c)

		var err error
		var questions []models.Question
		if user.Role == "admin" {
			questions, err = models.GetQuestions(appengine.NewContext(c.Request))
		} else {
			questions, err = models.GetQuestionsByUser(appengine.NewContext(c.Request), user.Id)
		}
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, questions)
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

		user := GetUserFromContext(c)

		if user.Role != "admin" && question.UserId != user.Id {
			c.String(http.StatusForbidden, "")
			return
		}

		c.JSON(http.StatusOK, question)
	})

	router.PUT("/:id", func(c *gin.Context) {
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

		user := GetUserFromContext(c)
		if user.Role != "admin" && question.UserId != user.Id {
			c.String(http.StatusForbidden, "")
			return
		}

		q, err := models.UpdateQuestion(appengine.NewContext(c.Request), id, c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, q)
	})

}
