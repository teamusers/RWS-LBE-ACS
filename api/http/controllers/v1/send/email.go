package send

import (
	"net/http"
	"rlp-email-service/api/http/requests"
	"rlp-email-service/api/http/responses"
	"rlp-email-service/codes"
	"rlp-email-service/utils"

	"github.com/gin-gonic/gin"
)

// Plaintext godoc.
// @Summary Send a plain text email
// @Description Sends a plain text email with the provided email, subject, plain text content, and optionally attachments.
// @Tags         send
// @Accept json
// @Produce json
// @Param request body requests.SendPlainTextEmailRequest true "Send email plain text body"
// @Success 200 {object} responses.ApiResponse[any] "Email sent"
// @Failure 400 {object} responses.ErrorResponse "Invalid request body"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Security     ApiKeyAuth
// @Router /api/v1/send/plain-text [post]
func PlainText(c *gin.Context) {
	var sendEmailRequest requests.SendPlainTextEmailRequest
	if err := c.ShouldBindJSON(&sendEmailRequest); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}
	err := utils.SendEmail(sendEmailRequest.Email, sendEmailRequest.Subject, sendEmailRequest.PlainText, false, sendEmailRequest.Attachments)
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

// EmailTemplate sends an email using a template.
// @Summary Send an email using a template
// @Description Sends an email using the specified template and data.
// @Tags         send
// @Accept json
// @Produce json
// @Param name path string true "Template name"
// @Param request body requests.SendTemplateEmailRequest true "Request body"
// @Success 200 {object} responses.ApiResponse[any] "Email sent"
// @Failure 400 {object} responses.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /send/template/{name} [post]
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

	err := utils.SendEmail(sendEmailRequest.Email, sendEmailRequest.Subject, content, true, sendEmailRequest.Attachments)
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
