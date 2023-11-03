package main

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type mortgageInfo struct {
	principal      float64
	interestRate   float64
	loanTermYears  int
	monthlyRate    float64
	monthlyPayment float64
	totalInterest  float64
	totalPayment   float64
}

func calculateMortgage(principal float64, interestRate float64, loanTermYears int) mortgageInfo {
	var info mortgageInfo

	monthlyRate := interestRate / 12.0 / 100.0
	numberOfPayments := float64(loanTermYears * 12)
	denominator := (1 + monthlyRate)

	powerFactor := math.Pow(denominator, numberOfPayments)
	monthlyPayment := principal * monthlyRate * powerFactor / (powerFactor - 1)
	totalPayment := monthlyPayment * numberOfPayments
	totalInterest := totalPayment - principal

	info.principal = principal
	info.interestRate = interestRate
	info.loanTermYears = loanTermYears
	info.monthlyRate = monthlyRate
	info.monthlyPayment = monthlyPayment
	info.totalInterest = totalInterest
	info.totalPayment = totalPayment

	return info
}

func renderScheduleInHtml(mInfo mortgageInfo) string {
	var scheduleHTML strings.Builder

	scheduleHTML.WriteString("<table><tr><th>Month</th><th>Interest</th><th>Principal</th><th>Remaining Balance</th></tr>")
	remainingBalance := mInfo.principal
	for month := 1; month <= mInfo.loanTermYears*12; month++ {
		monthlyInterest := remainingBalance * mInfo.monthlyRate
		monthlyPrincipal := mInfo.monthlyPayment - monthlyInterest
		remainingBalance -= monthlyPrincipal
		scheduleHTML.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>", month, monthlyInterest, monthlyPrincipal, remainingBalance))
	}
	scheduleHTML.WriteString("</table>")
	return scheduleHTML.String()

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
	mInfo := calculateMortgage(principal, interestRate, loanTermYears)

	// Open database
	db, err := sql.Open("sqlite3", "./hsl.db")

	if err != nil {
		panic("Cannot open database")
	}

	defer db.Close()

	// Store calculation in database
	// TODO

	// Respond to web client
	fmt.Fprintf(w, "Monthly Payment: %.2f, Total Payment: %.2f, Total Interest: %.2f", mInfo.monthlyPayment, mInfo.totalPayment, mInfo.totalInterest)
}

func main() {
	// Serve static files from the current directory
	http.HandleFunc("/", fileHandler)

	http.HandleFunc("/calculate", mortgageHandler)

	port := "8888"
	fmt.Println("Listening on port", port)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
