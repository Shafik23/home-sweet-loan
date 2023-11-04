package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type MonthData struct {
	Month            int     `json:"month"`
	MonthlyInterest  float64 `json:"monthlyInterest"`
	MonthlyPrincipal float64 `json:"monthlyPrincipal"`
	RemainingBalance float64 `json:"remainingBalance"`
}

type LoanRecord struct {
	Principal    float64 `json:"principal"`
	InterestRate float64 `json:"interestRate"`
	LoanTerm     int     `json:"loanTerm"`
}

type mortgageInfo struct {
	principal     float64
	interestRate  float64
	loanTermYears int
	monthlyRate   float64

	TotalInterest  float64     `json:"totalInterest"`
	MonthlyPayment float64     `json:"monthlyPayment"`
	TotalPayment   float64     `json:"totalPayment"`
	Schedule       []MonthData `json:"schedule"`
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
	info.MonthlyPayment = monthlyPayment
	info.TotalInterest = totalInterest
	info.TotalPayment = totalPayment

	// calculate the schedule
	remainingBalance := principal
	for month := 1; month <= loanTermYears*12; month++ {
		monthlyInterest := remainingBalance * info.monthlyRate
		monthlyPrincipal := info.MonthlyPayment - monthlyInterest
		remainingBalance -= monthlyPrincipal
		info.Schedule = append(info.Schedule, MonthData{month, monthlyInterest, monthlyPrincipal, remainingBalance})
	}

	return info
}

func jsonError(w http.ResponseWriter, message string, httpStatusCode int) {
	// Create a map for the error message
	errorMap := map[string]string{"error": message}
	// Convert the map to JSON
	jsonData, err := json.Marshal(errorMap)
	if err != nil {
		// If there is an error marshaling the JSON, fall back to http.Error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type header and the status code, then write the JSON error message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write(jsonData)
}

func mortgageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form values
	r.ParseForm()
	principalStr := r.FormValue("principal")
	interestRateStr := r.FormValue("interest_rate")
	loanTermYearsStr := r.FormValue("loan_term_years")

	// Convert form values to appropriate types
	principal, err := strconv.ParseFloat(principalStr, 64)
	if err != nil || principal < 0.1 || principal > 10000000 {
		jsonError(w, "Invalid principal format - please enter a number in the range [0.1, 10000000]", http.StatusBadRequest)
		return
	}

	interestRate, err := strconv.ParseFloat(interestRateStr, 64)
	if err != nil || interestRate < 0.1 || interestRate > 100 {
		jsonError(w, "Invalid interest rate - please enter a number in the range [0.1, 100]", http.StatusBadRequest)
		return
	}

	loanTermYears, err := strconv.Atoi(loanTermYearsStr)
	if err != nil || loanTermYears < 1 || loanTermYears > 100 {
		jsonError(w, "Invalid loan term - please enter a number in the range [1, 100]", http.StatusBadRequest)
		return
	}

	// Calculate mortgage
	mInfo := calculateMortgage(principal, interestRate, loanTermYears)

	// Store the user input in the DB
	err = storeLoan(principal, interestRate, loanTermYears)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("User-input inserted into DB successfully!")
	}

	jsonData, err := json.Marshal(mInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// This header actually triggers the history dropdown to reload!
	w.Header().Set("HX-Trigger", "reloadHistory")

	w.Write(jsonData)
}

func storeLoan(principal float64, interestRate float64, loanTermYears int) error {
	db, err := sql.Open("sqlite3", "./hsl.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Loans(Principal, InterestRate, LoanTerm) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(principal, interestRate, loanTermYears)
	return err
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./hsl.db")
	if err != nil {
		http.Error(w, "Failed to connect to the database.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query all loans from the database
	rows, err := db.Query("SELECT Principal, InterestRate, LoanTerm FROM Loans")
	if err != nil {
		http.Error(w, "Failed to fetch records.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var loans []LoanRecord

	// Iterate through all records and add them to the loans slice
	for rows.Next() {
		var loan LoanRecord

		err := rows.Scan(&loan.Principal, &loan.InterestRate, &loan.LoanTerm)

		if err != nil {
			http.Error(w, "Failed to read record.", http.StatusInternalServerError)
			return
		}

		loans = append(loans, loan)
	}

	// Check for any errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, "Failed during iteration of records.", http.StatusInternalServerError)
		return
	}

	// Convert loans slice to JSON
	jsonData, err := json.Marshal(loans)
	if err != nil {
		http.Error(w, "Failed to convert records to JSON.", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	// Serve static files from the current directory
	http.HandleFunc("/", fileHandler)

	http.HandleFunc("/calculate", mortgageHandler)
	http.HandleFunc("/fetchHistory", historyHandler)

	port := "8888"
	fmt.Printf("Listening on port %s - will exit if port is in use\n", port)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
