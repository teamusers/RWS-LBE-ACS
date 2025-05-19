package responses

type SMSResponse struct {
	Result struct {
		Status string `xml:"status"`
		Error  string `xml:"error"`
	} `xml:"result"`
}
