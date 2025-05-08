package utils

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func contentHash(body []byte) string {
	h := sha256.Sum256(body)
	return base64.StdEncoding.EncodeToString(h[:])
}

func computeSignature(stringToSign, base64Key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func SendEmail(email string, subject string, emailData string, isHtml bool) error {
	conn := os.Getenv("ACS_CONNECTION_STRING")
	var (
		endpoint, key string
	)
	for _, kv := range strings.Split(conn, ";") {
		if strings.HasPrefix(kv, "endpoint=") {
			endpoint = strings.TrimSuffix(strings.TrimPrefix(kv, "endpoint="), "/")
		}
		if strings.HasPrefix(kv, "accesskey=") {
			key = strings.TrimPrefix(kv, "accesskey=")
		}
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	host := u.Host
	pathAndQuery := "/emails:send?api-version=2023-03-31"
	content := map[string]string{
		"subject": subject,
	}
	if isHtml {
		content["html"] = emailData
	} else {
		content["plainText"] = emailData
	}

	bodyObj := map[string]interface{}{
		"senderAddress": os.Getenv("ACS_SENDER_ADDRESS"),
		"content":       content,
		"recipients": map[string][]map[string]string{
			"to": {{"address": email}},
		},
	}
	bodyBytes, _ := json.Marshal(bodyObj)

	gmt := time.FixedZone("GMT", 0)
	date := time.Now().In(gmt).Format(time.RFC1123)
	hash := contentHash(bodyBytes)

	stringToSign := fmt.Sprintf(
		"POST\n%s\n%s;%s;%s",
		pathAndQuery,
		date,
		host,
		hash,
	)

	sig, err := computeSignature(stringToSign, key)
	if err != nil {
		return err
	}
	authHeader := "HMAC-SHA256 SignedHeaders=x-ms-date;host;x-ms-content-sha256&Signature=" + sig

	req, err := http.NewRequestWithContext(context.Background(),
		"POST",
		endpoint+pathAndQuery,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return err
	}
	req.Host = host
	req.Header.Set("x-ms-date", date)
	req.Header.Set("x-ms-content-sha256", hash)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	io.ReadAll(resp.Body)
	return nil
}

func LoadTemplate(templateName string, data any) (string, error) {
	templatePath := fmt.Sprintf("templates/%s.html", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return tplBuffer.String(), nil
}
