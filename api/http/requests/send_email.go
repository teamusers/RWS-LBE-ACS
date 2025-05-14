package requests

type SendPlainTextEmailRequest struct {
	Email       string               `json:"email" binding:"required"`
	Subject     string               `json:"subject" binding:"required"`
	PlainText   string               `json:"plain_text" binding:"required"`
	Attachments *[]map[string]string `json:"attachments"`
}

type SendTemplateEmailRequest struct {
	Email       string                 `json:"email" binding:"required"`
	Subject     string                 `json:"subject" binding:"required"`
	Data        map[string]interface{} `json:"data" binding:"required"`
	Attachments *[]map[string]string   `json:"attachments"`
}
