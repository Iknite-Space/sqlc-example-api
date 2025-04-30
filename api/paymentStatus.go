package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type reference struct {
	Status string `json:"status"`
}

func CheckTransationStatus(token, ref string) string {

	resp, err := http.NewRequest(http.MethodGet, "https://demo.campay.net/api/transaction/"+ref+"/", nil)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

	//set the request headers
	resp.Header.Set("Authorization", "Token "+token)
	resp.Header.Set("Content-Type", "application/json")

	//sent the request using http.defaultclient
	client := &http.Client{}
	response, err := client.Do(resp)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	defer response.Body.Close()

	//read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//display the status code and text if not ok
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: %s", response.Status)
		log.Fatalf("Error: %s", string(body))
		return ""
	}

	//parse json data into struct
	var status reference
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatalf("Error parsing JSON Body %v", err)
	}
	//access the transaction reference
	return status.Status
}
