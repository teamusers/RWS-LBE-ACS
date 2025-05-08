package send

import (
	"net/http"
	"rlp-email-service/api/http/requests"
	"rlp-email-service/api/http/responses"
	"rlp-email-service/codes"
	"rlp-email-service/utils"

	"github.com/gin-gonic/gin"
)

func PlainText(c *gin.Context) {
	var sendEmailRequest requests.SendPlainTextEmailRequest
	if err := c.ShouldBindJSON(&sendEmailRequest); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}
	err := utils.SendEmail(sendEmailRequest.Email, sendEmailRequest.Subject, sendEmailRequest.PlainText, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
	resp := responses.ApiResponse[any]{
		Code:    codes.SUCCESSFUL,
		Message: "email sent",
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
func EmailTemplate(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	var sendEmailRequest requests.SendTemplateEmailRequest
	if err := c.ShouldBindJSON(&sendEmailRequest); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	content, loadTemplateErr := utils.LoadTemplate(name, sendEmailRequest.Data)
	if loadTemplateErr != nil {
		c.JSON(http.StatusInternalServerError, responses.UnsuccessfullLoadingTemplate())
		return
	}

	err := utils.SendEmail(sendEmailRequest.Email, sendEmailRequest.Subject, content, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
	resp := responses.ApiResponse[any]{
		Code:    codes.SUCCESSFUL,
		Message: "email sent",
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
