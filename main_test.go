package main

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateMortgage(t *testing.T) {
	tests := []struct {
		principal              float64
		interestRate           float64
		loanTermYears          int
		expectedMonthlyPayment float64
		expectedTotalPayment   float64
		expectedTotalInterest  float64
	}{
		{principal: 200000, interestRate: 3.5, loanTermYears: 30, expectedMonthlyPayment: 898.09, expectedTotalPayment: 323312.20, expectedTotalInterest: 123312.18},
		{principal: 150000, interestRate: 4.0, loanTermYears: 15, expectedMonthlyPayment: 1109.53, expectedTotalPayment: 199715.00, expectedTotalInterest: 49715.74},
	}

	// Tolerance of 1 dollar
	const tolerance = 1

	for _, test := range tests {
		mInfo := calculateMortgage(test.principal, test.interestRate, test.loanTermYears)
		if !closeEnough(mInfo.MonthlyPayment, test.expectedMonthlyPayment, tolerance) {
			t.Errorf("Expected monthly payment of %.2f but got %.2f", test.expectedMonthlyPayment, mInfo.MonthlyPayment)
		}
		if !closeEnough(mInfo.TotalPayment, test.expectedTotalPayment, tolerance) {
			t.Errorf("Expected total payment of %.2f but got %.2f", test.expectedTotalPayment, mInfo.TotalPayment)
		}
		if !closeEnough(mInfo.TotalInterest, test.expectedTotalInterest, tolerance) {
			t.Errorf("Expected total interest of %.2f but got %.2f", test.expectedTotalInterest, mInfo.TotalInterest)
		}
	}
}

func closeEnough(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestInputValidation(t *testing.T) {
	tests := []struct {
		input          string
		expectedStatus int
		expectedBody   string
	}{
		{input: "principal=-100&interest_rate=4.5&loan_term_years=30", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid principal format - please enter a number in the range [0.1, 10000000]"}`},
		{input: "principal=200000&interest_rate=-4.5&loan_term_years=30", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid interest rate - please enter a number in the range [0.1, 100]"}`},
		{input: "principal=200000&interest_rate=4.5&loan_term_years=-30", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid loan term - please enter a number in the range [1, 100]"}`},

		{input: "principal=200000&interest_rate=4.5&loan_term_years=absdf", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid loan term - please enter a number in the range [1, 100]"}`},
		{input: "principal=200000&interest_rate=dfkjskd&loan_term_years=30", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid interest rate - please enter a number in the range [0.1, 100]"}`},
		{input: "principal=sdkfjsd&interest_rate=4.5&loan_term_years=30", expectedStatus: http.StatusBadRequest, expectedBody: `{"error":"Invalid principal format - please enter a number in the range [0.1, 10000000]"}`},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/calculate", strings.NewReader(test.input))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mortgageHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, test.expectedStatus)
		}

		if rr.Body.String() != test.expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), test.expectedBody)
		}
	}
}
