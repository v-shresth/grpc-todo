package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"todo-grpc/utils"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type EmailClient struct {
	apiKey        string
	client        HTTPClient
	senderAddress string
}

type sendEmail struct {
	Key     string        `json:"key"`
	Message *emailMessage `json:"message"`
}

type emailMessage struct {
	Html      string       `json:"html"`
	Subject   string       `json:"subject"`
	FromEmail string       `json:"from_email"`
	To        []*recipient `json:"to"`
}

type recipient struct {
	Email string `json:"email"`
}

const base_url = "https://mandrillapp.com/api/1.0"

func NewEmailClient(config utils.EnvConfig) *EmailClient {
	return &EmailClient{
		apiKey:        config.GetMailChimpApiKey(),
		client:        &http.Client{},
		senderAddress: config.GetSenderEmailAddress(),
	}
}

func (m *EmailClient) SendEmail(
	receiverEmail string,
	emailHtmlString string,
	subject string,
) error {
	data := &sendEmail{
		Key: m.apiKey, Message: &emailMessage{
			Html:      emailHtmlString,
			Subject:   subject,
			FromEmail: m.senderAddress,
			To: []*recipient{
				{Email: receiverEmail},
			},
		},
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		base_url+"/messages/send",
		bytes.NewBuffer(dataBytes),
	)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer func() {
		if req.Body != nil {
			if err := req.Body.Close(); err != nil {
				fmt.Println("error closing request body")
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error with mail chimp request: %s", resp.Status)
	}

	return nil
}
