package requests

type SendEmailRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Subject   string `json:"subject" binding:"required,subject"`
	PlainText string `json:"plain_text" binding:"required,plain_text"`
	Html      string `json:"html" binding:"required,html"`
}
