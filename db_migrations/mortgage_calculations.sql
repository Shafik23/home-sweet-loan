CREATE TABLE mortgage_calculations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    principal_amount FLOAT,
    interest_rate FLOAT,
    loan_term_years INTEGER,
    monthly_payment FLOAT,
    total_payment FLOAT,
    total_interest FLOAT
);

