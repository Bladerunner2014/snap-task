package main


import (
    "bytes"
    "encoding/json"
    "net/http"
    // "net/http/httptest"
    "testing"
    // "time"

 
)

func TestJobAPI(t *testing.T) {
	type JobRequest struct {
		URL     string            `json:"url"`
		Pattern map[string]string `json:"pattern"`
		Interval int              `json:"interval"`
	}
	
	// Create the request payload
	payload := JobRequest{
		URL: "http://api.example.com/data",
		Pattern: map[string]string{
			"name":   "string",
			"age":    "float64",
			"height": "float64",
		},
		Interval: 5,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	// Send the POST request to the API
	resp, err := http.Post("http://0.0.0.0:8080/job", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, but got %d", resp.StatusCode)
	}
}
