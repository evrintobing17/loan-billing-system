package models

import "time"

type Payment struct {
	ID             int       `json:"id"`
	LoanID         int       `json:"loan_id"`
	Amount         float64   `json:"amount"`
	PaymentDate    time.Time `json:"payment_date"`
	IdempotencyKey string    `json:"idempotency_key"`
}

type PaymentInstallment struct {
	PaymentID     int `json:"payment_id"`
	InstallmentID int `json:"installment_id"`
}

type PaymentRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
