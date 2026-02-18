package loan

import (
	"context"

	"github.com/evrintobing17/loan-billing-system/models"
)

type LoanRepository interface {
	Create(ctx context.Context, loan *models.Loan, installments []models.Installment) error
	GetByID(ctx context.Context, id int) (*models.Loan, error)
	GetInstallments(ctx context.Context, loanID int) ([]models.Installment, error)
	UpdateInstallmentsPaid(ctx context.Context, installmentIDs []int) error
}
