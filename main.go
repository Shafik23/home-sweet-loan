package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

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

// Calculates mortgage related information given the principal amount, interest rate, and loan term in years
// Returns calculated data including monthly payment, total payment, total interest, and a payment schedule
func calculateMortgage(principal float64, interestRate float64, loanTermYears int) mortgageInfo {
	var info = mortgageInfo{
		principal:     principal,
		interestRate:  interestRate,
		loanTermYears: loanTermYears,
	}

	info.monthlyRate = interestRate / 12.0 / 100.0
	numberOfPayments := float64(loanTermYears * 12)
	denominator := 1 + info.monthlyRate
	powerFactor := math.Pow(denominator, numberOfPayments)

	// Equation for mortgage payments, from pdf
	info.MonthlyPayment = principal * info.monthlyRate * powerFactor / (powerFactor - 1)

	// Calculate total payment and total interest
	info.TotalPayment = info.MonthlyPayment * numberOfPayments
	info.TotalInterest = info.TotalPayment - principal

	// Calculate the payment schedule
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

func jsonResponse(w http.ResponseWriter, data interface{}, httpStatusCode int) {
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Fallback to http.Error in case of marshaling error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type header and the status code, then write the JSON data
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
		// Note: a unique constraint violation will cause an error here,
		// even though it is expected behavior. This is not a problem we
		// need to report to the user.
		fmt.Println(err)
	} else {
		fmt.Println("User-input inserted into DB successfully!")
	}

	// This htmx header triggers the history dropdown to reload!
	w.Header().Set("HX-Trigger", "reloadHistory")

	jsonResponse(w, mInfo, http.StatusOK)
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
	rows, err := db.Query("SELECT Principal, InterestRate, LoanTerm FROM Loans ORDER BY CreatedAt DESC LIMIT 100")
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

	jsonResponse(w, loans, http.StatusOK)
}

type MarketRateResponse struct {
	Rate string `json:"rate"`
}

func marketRateHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := fetchCurrentMarketMortgageRate()
	if err != nil {
		// You can use the jsonError function if it's appropriate
		jsonError(w, "Failed to fetch the market mortgage rate", http.StatusInternalServerError)
		return
	}

	// Create response object
	response := MarketRateResponse{Rate: rate}

	jsonResponse(w, response, http.StatusOK)
}

func fetchCurrentMarketMortgageRate() (string, error) {
	// URL of the market mortgage rate data
	url := "https://fred.stlouisfed.org/data/MORTGAGE30US.txt"

	// HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert body to string and split by newline
	lines := strings.Split(string(body), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("unexpected format of the data")
	}

	// The last line contains the latest rate
	lastLine := lines[len(lines)-2] // -2 because the last element is empty due to newline at the end
	parts := strings.Fields(lastLine)
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected format of the last line")
	}

	// Return the last part which should be the rate
	return parts[len(parts)-1], nil
}

func main() {
	// Serve static files from the current directory
	http.HandleFunc("/", fileHandler)

	http.HandleFunc("/calculate", mortgageHandler)
	http.HandleFunc("/fetchHistory", historyHandler)
	http.HandleFunc("/currentMarketRate", marketRateHandler)

	port := "8888"
	fmt.Printf("Listening on port %s - will abort if port is in use.\n", port)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
