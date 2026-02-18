package usecase

import (
	"context"
	"time"

	"github.com/evrintobing17/loan-billing-system/internal/loan"
	"github.com/evrintobing17/loan-billing-system/models"
)

type loanUseCase struct {
	loanRepo loan.LoanRepository
}

func NewLoanUseCase(loanRepo loan.LoanRepository) loan.LoanUsecase {
	return &loanUseCase{loanRepo: loanRepo}
}

func (uc *loanUseCase) CreateLoan(ctx context.Context, principal, rate float64, termWeeks int, startDate time.Time) (*models.Loan, error) {
	weeklyAmount := calculateWeeklyAmount(principal, rate, termWeeks)

	loan := &models.Loan{
		Principal:    principal,
		InterestRate: rate,
		TermWeeks:    termWeeks,
		WeeklyAmount: weeklyAmount,
		StartDate:    startDate,
		IsActive:     true,
	}
	installments := generateInstallments(loan)
	err := uc.loanRepo.Create(ctx, loan, installments)

	return loan, err
}

func (uc *loanUseCase) GetOutstanding(ctx context.Context, loanID int) (float64, error) {
	installments, err := uc.loanRepo.GetInstallments(ctx, loanID)
	if err != nil {
		return 0, err
	}
	var outstanding float64
	for _, inst := range installments {
		if !inst.Paid {
			outstanding += inst.Amount
		}
	}
	return outstanding, nil
}

func (uc *loanUseCase) IsDelinquent(ctx context.Context, loanID int) (bool, error) {
	installments, err := uc.loanRepo.GetInstallments(ctx, loanID)
	if err != nil {
		return false, err
	}
	today := time.Now().Truncate(24 * time.Hour)
	var unpaidPastDue []models.Installment
	for _, inst := range installments {
		if !inst.Paid && !inst.DueDate.After(today) {
			unpaidPastDue = append(unpaidPastDue, inst)
		}
	}
	// check for consecutive weeks
	for i := 1; i < len(unpaidPastDue); i++ {
		if unpaidPastDue[i].WeekNumber == unpaidPastDue[i-1].WeekNumber+1 {
			return true, nil
		}
	}
	return false, nil
}

// Helper functions
func calculateWeeklyAmount(principal, rate float64, weeks int) float64 {
	interest := principal * rate / 100
	total := principal + interest
	return total / float64(weeks)
}

func generateInstallments(loan *models.Loan) []models.Installment {
	installments := make([]models.Installment, loan.TermWeeks)
	for i := 0; i < loan.TermWeeks; i++ {
		dueDate := loan.StartDate.AddDate(0, 0, (i+1)*7)
		installments[i] = models.Installment{
			WeekNumber: i + 1,
			DueDate:    dueDate,
			Amount:     loan.WeeklyAmount,
			Paid:       false,
		}
	}
	return installments
}
