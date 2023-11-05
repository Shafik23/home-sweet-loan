CREATE TABLE Loans (
    Principal REAL NOT NULL,
    InterestRate REAL NOT NULL,
    LoanTerm INTEGER NOT NULL,
    CreatedAt datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (Principal, InterestRate, LoanTerm)
);