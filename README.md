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




## Architecture Design Record
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
  - A secondary reason was *exploration*! I had not had any experience with it and used this assignment as an opportunity to learn a new framework.
  - In retrospect, it was probably not the best decision. A more appropriate choice would have been something like **Svelte**.

### Error Handling

Input validation is performed on the backend. Errors are returned as JSON objects with a descriptive error message, to be rendered to the screen by the UI.

### Testing

Unit tests have been written for the calculateMortgage function to ensure accurate mortgage calculations.
Additional tests have been written for the mortgageHandler function to ensure proper input validation and error handling.
