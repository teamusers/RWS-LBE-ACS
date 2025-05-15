package send

import (
	"fmt"
	"net/http"
	"rlp-email-service/api/http/requests"
	"rlp-email-service/api/http/responses"
	"rlp-email-service/api/http/services"
	"rlp-email-service/codes"

	"github.com/gin-gonic/gin"
)

// Sms godoc
// @Summary      Sends sms
// @Description  Sends sms to the given phone numbers
// @Tags         send
// @Accept       json
// @Produce      json
// @Param        request  body      requests.SendSMSRequest          true  "Send SMS request body"
// @Success      200      {object}  responses.ApiResponse[any] 					 "nil"
// @Failure      400      {object}  responses.ErrorResponse            "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                             "Unauthorized – API key missing or invalid"
// @Failure      500      {object}  responses.ErrorResponse                     "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/v1/send/sms [post]
func SMS(c *gin.Context) {
	var req requests.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}
	err := services.SendSMS(req.Receivers, req.Content)
	if err != nil {
		fmt.Println("SendSMS error:", err)
		c.JSON(http.StatusInternalServerError, responses.DefaultResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	resp := responses.ApiResponse[any]{
		Code:    codes.SUCCESSFUL,
		Message: "sms sent",
		Data:    nil,
	}
	c.JSON(http.StatusOK, resp)
}
