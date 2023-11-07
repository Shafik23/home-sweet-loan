package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestCalculateEndpoint(t *testing.T) {
	testCases := []struct {
		principal      string
		interestRate   string
		loanTermYears  string
		expectedStatus int
		description    string
	}{
		{"200000", "3.5", "30", http.StatusOK, "valid input"},
		{"-2000", "3.5", "30", http.StatusBadRequest, "negative principal"},
		{"200000", "-3.5", "30", http.StatusBadRequest, "negative interest rate"},
		{"200000", "3.5", "-30", http.StatusBadRequest, "negative loan term"},
		{"200000", "3.5", "", http.StatusBadRequest, "missing loan term"},
		{"", "3.5", "30", http.StatusBadRequest, "missing principal"},
		{"200000", "", "30", http.StatusBadRequest, "missing interest rate"},
		{"notanumber", "3.5", "30", http.StatusBadRequest, "non-numeric principal"},
		{"200000", "notanumber", "30", http.StatusBadRequest, "non-numeric interest rate"},
		{"200000", "3.5", "notanumber", http.StatusBadRequest, "non-numeric loan term"},
		{"0", "3.5", "30", http.StatusBadRequest, "zero principal"},
		{"200000", "0", "30", http.StatusBadRequest, "zero interest rate"},
		{"200000", "3.5", "0", http.StatusBadRequest, "zero loan term"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8888/calculate?principal=%s&interest_rate=%s&loan_term_years=%s",
				url.QueryEscape(tc.principal), url.QueryEscape(tc.interestRate), url.QueryEscape(tc.loanTermYears))
			resp, err := http.Get(url)

			if err != nil {
				t.Fatalf("Failed to send request for '%s': %v", tc.description, err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("For '%s', expected status %v but got %v", tc.description, tc.expectedStatus, resp.Status)
			}
		})
	}
}

func TestFetchHistoryEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8888/fetchHistory")

	if err != nil {
		t.Fatalf("Failed to send request to fetch history: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK but got %v when fetching history", resp.Status)
	}

	// Additional tests for fetchHistory can be added here.
}
