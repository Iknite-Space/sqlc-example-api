package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const campayUrl = "https://demo.campay.net/api/collect/"

type paymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	From        string `json:"from"`
	Description string `json:"description"`
}

type paymentResponse struct {
	Reference string `json:"reference"`
}

func requestPayment(phone, amount, description, token string) (string, error) {
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

	req, err := http.NewRequest("POST", campayUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var paymentResp paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to initiate payment: %s", resp.Status)
	}

	return paymentResp.Reference, nil
}
