CREATE TABLE Loans (
    Principal REAL NOT NULL,
    InterestRate REAL NOT NULL,
    LoanTerm INTEGER NOT NULL,
    PRIMARY KEY (Principal, InterestRate, LoanTerm)
);