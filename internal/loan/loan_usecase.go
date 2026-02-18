package loan

import (
	"context"
	"time"

	"github.com/evrintobing17/loan-billing-system/models"
)

type LoanUsecase interface {
	CreateLoan(ctx context.Context, principal, rate float64, termWeeks int, startDate time.Time) (*models.Loan, error)
	GetOutstanding(ctx context.Context, loanID int) (float64, error)
	IsDelinquent(ctx context.Context, loanID int) (bool, error)
}
