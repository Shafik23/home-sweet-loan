package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCalculateEndpointJSONShape(t *testing.T) {
	url := "http://localhost:8888/calculate?principal=200000&interest_rate=3.5&loan_term_years=30"
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	var info mortgageInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	require.NoError(t, err)

	assert.NotZero(t, info.MonthlyPayment)
	assert.NotZero(t, info.TotalPayment)
	assert.NotZero(t, info.TotalInterest)
	assert.NotEmpty(t, info.Schedule)
}

func TestFetchHistoryEndpointJSONShape(t *testing.T) {
	resp, err := http.Get("http://localhost:8888/fetchHistory")
	require.NoError(t, err)
	defer resp.Body.Close()

	var loans []LoanRecord
	err = json.NewDecoder(resp.Body).Decode(&loans)
	require.NoError(t, err)

	for _, loan := range loans {
		assert.NotZero(t, loan.Principal)
		assert.NotZero(t, loan.InterestRate)
		assert.NotZero(t, loan.LoanTerm)
	}
}

func TestMarketRateEndpointJSONShape(t *testing.T) {
	resp, err := http.Get("http://localhost:8888/currentMarketRate")
	require.NoError(t, err)
	defer resp.Body.Close()

	var marketRate MarketRateResponse
	err = json.NewDecoder(resp.Body).Decode(&marketRate)
	require.NoError(t, err)

	assert.NotEmpty(t, marketRate.Rate)
}
