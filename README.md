# Home Sweet Loan
Simple Home Loan Schedule Calculator!

### Dependencies

- SQLite3
- Golang v1.20
- Htmx
- Water.css

### Setup
- Set up DB schema by running `db_migrations/loans.sql`.
  - Run `sqlite3 hsl.db` (name must be exact).
  - Run `.read ./db_migrations/loans.sql`
- Build: `go build`
- Run: `./home-sweet-loan`

### Development
- Pretty straightforward:
  - Backend entry-point: main.go
  - Frontend entry-point: hsl.html
  - DB: hsl.db
- Sorry, no hot-reload or live-reload!


## Architecture 
This is a simple 3-tier web application that computes the amortization schedule for a loan given a few user inputs. 

It is built using **Go** for the backend, **SQLite** for the database, and **HTMX** for the frontend. The goal is to keep the implementation very simple.

### Components

- **Database:** SQLite, a lightweight disk-based database that doesn’t require a separate server process.
	- Quite a mature technology that focuses on being lightweight and simple to deploy and manage.

- **Backend API:** Written in Go, handles HTTP requests, performs mortgage calculations, and interacts with the SQLite database.
  - Go was chosen for its simple syntax, excellent support for web-backends, and ease of deployment to target environments.
  - There are 2 endpoints - both return data in JSON format:
    - `/calculate`: Does the heavy lifting of computing the loan-schedule along with other information.
    - `/fetchHistory`: Fetches user-input history stored in the DB.

- **Frontend:** Implemented using HTMX, quite a minimal JavaScript library for accessing HTML templates and updating the DOM.
  - This was chosen in alignment with the goal of keeping everything very simple.
  - A secondary reason was *exploration*! I used this assignment as an opportunity to learn a new framework.
  - In retrospect, it was probably not the best decision. A more appropriate choice would have been something like **Svelte**.

### Error Handling

Input validation is performed on the backend. Errors are returned as JSON objects with a descriptive error message, to be rendered on the screen by the UI.

### Testing

Testing is divided into 3 categories: to run all tests, use the command `go test -v`.

- Unit Testing:
  - `main_test.go`
  - These test the main "business logic" of the application.
  - Given more time, more tests can be written to test each function in the backend. Unit tests can also be expanded to include UI components/etc.
- API Testing:
  - `api_test.go`
  - These test the backend API.
  - Given more time, this can be vastly expanded to include many more edge cases/etc.
- End-to-end Testing:
  - `e2e_test.go`
  - These theoretically test the entire app, end to end (i.e. User Interface -> Database, and back).
  - In the real world we would use a robust framework for doing this like Selenium or my favorite, PlayWright.

## User Stories

### Feature: Mortgage Calculation
User Story: 

As a user, I want to calculate the monthly payment of a mortgage so that I can understand if I can afford the mortgage.

Acceptance Criteria:

- The calculator should allow the user to input the loan amount, interest rate, and loan term in years.
- The calculator must validate the input to ensure it makes sense mathematically.
- The calculator should display a clear error message if the input is invalid.
- The monthly payment should be calculated and displayed to the user upon submission of valid data.

### Feature: Amortization Schedule
User Story: 

As a user, I want to see an amortization schedule so that I can know how much principal and interest I will pay each month.

Acceptance Criteria:

- The schedule should be displayed after the monthly payment is calculated.
- The schedule must list all payments for the entire term of the loan.
Each entry in the amortization schedule should show the payment number, interest amount, principal amount, and remaining balance.
- The final payment should reduce the balance to zero.

### Feature: Historical View of Calculations
User Story: 

As a user, I want to see a history of my past calculations so that I can track changes and compare different mortgage options.

Acceptance Criteria:

- The user should be able to access a history view from the main interface.
- The user should be able to select a past calculation to view its details.

## Potential Concerns
Obviously this is not a serious application and there are many potential pitfalls and deficiencies in it. Here are some of them (in no particular order):
  - Lack of robust e2e testing.
  - Not using HTMX in the intended idiomatic way (rendering html on the backend).
  - No explicit CORS policy.
  - Implicit coupling between frontend and backend:
    - The frontend has an implicit dependency on the structure of the server's JSON response. If the response structure changes, the JS code may fail.
    - The regex used in the JS for extracting values from the dropdown could fail if the formatting of the option text changes. This coupling requires that the format of the dropdown text be strictly maintained.
  - CSS hardcoded dimensions.
    - This may make it difficult for the UI to render correctly on different screen sizes.
  - Very basic/inefficient Database interaction.
  - No CSRF protection.
  - Minimal input sanitization.
  - Lack of ARIA attributes.