package campay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
	Client  *http.Client //the actual HTTP client used to make requests
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		Client:  &http.Client{},
	}
}

type paymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	From        string `json:"from"`
	Description string `json:"description"`
}

type paymentResponse struct {
	Reference string `json:"reference"`
}

type reference struct {
	Status string `json:"status"`
}

// RequestPayment() is a method that uses the http.Client tool stored inside the Client struct to make the actual request
func (c *Client) RequestPayment(phone, amount, description string) (string, error) {
	paymentReq := paymentRequest{
		Amount:      amount,
		Currency:    "XAF",
		From:        phone,
		Description: description,
	}

	reqBody, err := json.Marshal(paymentReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/collect/", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to initiate payment: %s", resp.Status)
	}

	var paymentResp paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return "", err
	}

	return paymentResp.Reference, nil
}

func (c *Client) CheckTransactionStatus(ref string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/transaction/"+ref+"/", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s\n%s", resp.Status, string(body))
	}

	var status reference
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", err
	}
	return status.Status, nil
}
