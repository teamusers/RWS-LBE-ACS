package requests

type SendSMSRequest struct {
	Receivers []int  `json:"receivers" binding:"required" example:"6581234567,6581234568"`
	Content   string `json:"content" binding:"required" example:"Your verification code is 123456"`
}
