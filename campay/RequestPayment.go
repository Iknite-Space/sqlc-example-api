package campay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PaymentRequest struct {
	From               string `json:"from"`
	Amount             string `json:"amount"`
	Description        string `json:"description"`
	External_Reference string `json:"external_reference"`
}

type PaymentResponse struct {
	Reference string `json:"reference"`
	Ussd_Code string `json:"ussd_code"`
}

func Payment(apikey string, number string, amount string, description string, external_reference string) PaymentResponse {

	number = "237" + number

	client := &http.Client{}

	iraq := PaymentRequest{
		From:               number,
		Amount:             amount,
		Description:        description,
		External_Reference: external_reference,
	}

	postbody, _ := json.Marshal(iraq)

	postbodyparam := bytes.NewBuffer(postbody)

	req, err := http.NewRequest("POST", "https://demo.campay.net/api/collect/", postbodyparam)

	if err != nil {
		fmt.Printf(" could not make a new request")
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apikey))
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Invalid Request, check POST request credentials")
		log.Fatal(err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var makepay PaymentResponse
	if err := json.NewDecoder(response.Body).Decode(&makepay); err != nil {
		log.Printf("failed to decode response: %v", err)
	}

	return makepay
	// defer response.Body.Close()

	// var makepay PaymentResponse
	// json.NewDecoder(response.Body).Decode(&makepay)
	// return makepay
}
