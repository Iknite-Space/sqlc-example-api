package campay

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PaymentStatus struct {
	Reference string `json:"reference"`
	Status    string `json:"status"`
}

func Status(apikey string, reference string) PaymentStatus {
	client := &http.Client{}
	var url = fmt.Sprintf("https://demo.campay.net/api/transaction/%s/", reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf(" could not make a new request")
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", apikey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf(" check get request method")
		log.Fatal(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var state PaymentStatus
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		log.Printf("failed to decode response: %v", err)

	}
	return state
	// defer resp.Body.Close()

	// var state PaymentStatus
	// json.NewDecoder(resp.Body).Decode(&state)
	// return state

}
