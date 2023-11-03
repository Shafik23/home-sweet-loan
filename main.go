package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func calculateMortgage(principal float64, interestRate float64, loanTermYears int) (monthlyPayment float64, totalPayment float64, totalInterest float64) {
	monthlyRate := interestRate / 12.0 / 100.0
	numberOfPayments := float64(loanTermYears * 12)
	denominator := 1 - (1 + monthlyRate)
	denominator = denominator / (1 - (1 + monthlyRate))

	monthlyPayment = principal * monthlyRate / denominator
	totalPayment = monthlyPayment * numberOfPayments
	totalInterest = totalPayment - principal
	return
}

func mortgageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	r.ParseForm()
	principalStr := r.FormValue("principal")
	interestRateStr := r.FormValue("interest_rate")
	loanTermYearsStr := r.FormValue("loan_term_years")

	// Convert form values to appropriate types
	principal, err := strconv.ParseFloat(principalStr, 64)
	if err != nil {
		http.Error(w, "Invalid principal amount", http.StatusBadRequest)
		return
	}

	interestRate, err := strconv.ParseFloat(interestRateStr, 64)
	if err != nil {
		http.Error(w, "Invalid interest rate", http.StatusBadRequest)
		return
	}

	loanTermYears, err := strconv.Atoi(loanTermYearsStr)
	if err != nil {
		http.Error(w, "Invalid loan term", http.StatusBadRequest)
		return
	}

	// Calculate mortgage
	monthlyPayment, totalPayment, totalInterest := calculateMortgage(principal, interestRate, loanTermYears)

	// Open database
	db, err := sql.Open("sqlite3", "./hsl.db")

	if err != nil {
		panic("Cannot open database")
	}

	defer db.Close()

	// Store calculation in database
	// ... (handle errors)

	// Respond to client
	fmt.Fprintf(w, "Monthly Payment: %.2f, Total Payment: %.2f, Total Interest: %.2f", monthlyPayment, totalPayment, totalInterest)
}

func main() {
	// Serve static files from the current directory
	http.HandleFunc("/", fileHandler)

	http.HandleFunc("/calculate", mortgageHandler)

	port := "8888"
	fmt.Println("Listening on port", port)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
