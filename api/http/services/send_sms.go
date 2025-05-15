package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"rlp-email-service/api/http/responses"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL    = "https://www.smsdome.com"
	sendSMSUrl = "/api/http/sendsms.aspx"
)

func SendSMS(receivers []int, content string) error {
	receiverStrs := make([]string, len(receivers))
	for i, r := range receivers {
		receiverStrs[i] = strconv.Itoa(r)
	}
	joinedReceivers := strings.Join(receiverStrs, ",")

	queryParams := map[string]string{
		"receivers": joinedReceivers,
		"content":   content,
	}

	resp, err := buildSendDomeHttpClient("POST", sendSMSUrl, queryParams, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var smsResp responses.SMSResponse
	if err := xml.Unmarshal(body, &smsResp); err != nil {
		return err
	}

	if strings.ToUpper(smsResp.Result.Status) == "NOK" {
		return errors.New(smsResp.Result.Error)
	}

	return nil
}

func buildSendDomeHttpClient(httpMethod string, path string, queryParams map[string]string, payload any) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	fullURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, path))
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	query := fullURL.Query()
	query.Set("appid", os.Getenv("SENDBOX_APPID"))
	query.Set("appsecret", os.Getenv("SENDBOX_SECRET"))

	for param, value := range queryParams {
		query.Set(param, value)
	}
	fullURL.RawQuery = query.Encode()

	var req *http.Request
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling payload: %w", err)
		}
		req, err = http.NewRequest(httpMethod, fullURL.String(), bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(httpMethod, fullURL.String(), nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := client.Do(req)
	return resp, err
}
