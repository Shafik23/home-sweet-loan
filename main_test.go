package main

import (
	"math"
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

func round(x float64, prec int) float64 {
	pow10 := math.Pow(10, float64(prec))
	return math.Round(x*pow10) / pow10
}

func closeEnough(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
