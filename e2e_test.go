package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHomePageHasCalculateButton(t *testing.T) {
	resp, err := http.Get("http://localhost:8888")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK but got %v", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Convert the body to a string for easier searching
	bodyString := string(body)

	// Check for the presence of the Calculate button
	if !strings.Contains(bodyString, `<button type="submit">Calculate</button>`) {
		t.Errorf("Did not find the Calculate button on the home page")
	}
}
