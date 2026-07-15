package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailjetClient struct {
	apiKey      string
	secretKey   string
	senderEmail string
	senderName  string
	baseURL     string
	httpClient  *http.Client
}

func NewMailjetClient(apiKey, secretKey, senderEmail, senderName, baseURL string) *MailjetClient {
	return &MailjetClient{
		apiKey:      apiKey,
		secretKey:   secretKey,
		senderEmail: senderEmail,
		senderName:  senderName,
		baseURL:     baseURL,
		httpClient:  &http.Client{},
	}
}

type mailjetRecipient struct {
	Email string `json:"Email"`
	Name  string `json:"Name,omitempty"`
}

type mailjetMessage struct {
	From     mailjetRecipient   `json:"From"`
	To       []mailjetRecipient `json:"To"`
	Subject  string             `json:"Subject"`
	TextPart string             `json:"TextPart"`
}

type mailjetRequest struct {
	Messages []mailjetMessage `json:"Messages"`
}

func (c *MailjetClient) SendEmail(toEmail, toName, subject, textContent string) error {
	reqBody := mailjetRequest{
		Messages: []mailjetMessage{
			{
				From:     mailjetRecipient{Email: c.senderEmail, Name: c.senderName},
				To:       []mailjetRecipient{{Email: toEmail, Name: toName}},
				Subject:  subject,
				TextPart: textContent,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal mailjet request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build mailjet request: %w", err)
	}
	req.SetBasicAuth(c.apiKey, c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send mailjet request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mailjet returned status %d", resp.StatusCode)
	}

	return nil
}
