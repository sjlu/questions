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
		// categoryId := c.Request.Form.Get("category_id")

		user := GetUserFromContext(c)
		userId := user.Id
		if user.Role == "admin" {
			userId = 0
		}

		var err error
		var questions []models.Question
		questions, err = models.GetQuestions(appengine.NewContext(c.Request), userId)
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
			c.AbortWithStatus(http.StatusForbidden)
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
			c.AbortWithStatus(http.StatusForbidden)
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
