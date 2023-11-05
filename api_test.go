package main

import (
	"net/http"
	"testing"
)

func TestCalculate(t *testing.T) {
	resp, err := http.Get("http://localhost:8888/calculate?principal=200000&interest_rate=3.5&loan_term_years=30")

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK but got %v", resp.Status)
	}
}