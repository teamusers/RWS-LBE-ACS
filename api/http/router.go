package http

import (
	v1 "rlp-email-service/api/http/controllers/v1"
	"rlp-email-service/api/http/controllers/v1/send"
	"rlp-email-service/api/interceptor"

	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {

	v1Group := e.Group("/v1")
	v1Group.POST("/auth", v1.AuthHandler)
	sendGroup := v1Group.Group("/send", interceptor.HttpInterceptor())
	{
		sendGroup.POST("/plain-text", send.PlainText)
		sendGroup.POST("/template/:name", send.EmailTemplate)
	}
}
