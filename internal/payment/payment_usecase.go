package payment

import "context"

type PaymentUsecase interface {
	MakePayment(ctx context.Context, loanID int, amount float64, idempotencyKey string) error
}
