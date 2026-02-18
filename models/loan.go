package models

import "time"

type Loan struct {
	ID           int       `json:"id"`
	Principal    float64   `json:"principal"`
	InterestRate float64   `json:"interest_rate"`
	TermWeeks    int       `json:"term_weeks"`
	WeeklyAmount float64   `json:"weekly_amount"`
	StartDate    time.Time `json:"start_date"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type Installment struct {
	ID         int       `json:"id"`
	LoanID     int       `json:"loan_id"`
	WeekNumber int       `json:"week_number"`
	DueDate    time.Time `json:"due_date"`
	Amount     float64   `json:"amount"`
	Paid       bool      `json:"paid"`
}

type CreateLoanRequest struct {
	Principal    float64 `json:"principal" binding:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" binding:"required,gt=0"`
	TermWeeks    int     `json:"term_weeks" binding:"required,gt=0"`
	StartDate    string  `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
}
