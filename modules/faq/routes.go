package faq

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	faq := app.Group("/faq")
	faq.Use(middleware.AuthMiddleware())
	{
		faq.GET("/list", GetFaqList)
	}
}
