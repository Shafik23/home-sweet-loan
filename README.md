# home-sweet-loan
Simple Home Loan Calculator!

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
     