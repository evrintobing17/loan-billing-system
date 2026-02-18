package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/evrintobing17/loan-billing-system/internal/loan"
	"github.com/evrintobing17/loan-billing-system/internal/payment"
	"github.com/evrintobing17/loan-billing-system/models"
	"github.com/evrintobing17/loan-billing-system/pkg/idempotency"
)

type paymentUseCase struct {
	paymentRepo payment.PaymentRepository
	loanRepo    loan.LoanRepository
	idempStore  idempotency.Store
}

func NewPaymentUseCase(pr payment.PaymentRepository, lr loan.LoanRepository, idemp idempotency.Store) payment.PaymentUsecase {
	return &paymentUseCase{
		paymentRepo: pr,
		loanRepo:    lr,
		idempStore:  idemp,
	}
}

func (uc *paymentUseCase) MakePayment(ctx context.Context, loanID int, amount float64, idempotencyKey string) error {
	// Idempotency check
	if idempotencyKey != "" {
		exists, err := uc.idempStore.Exists(ctx, idempotencyKey)
		if err != nil {
			return err
		}
		if exists {
			// Already processed, silently succeed
			return nil
		}
	}

	// Validate loan
	loan, err := uc.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return err
	}
	if amount <= 0 || int(amount*100)%int(loan.WeeklyAmount*100) != 0 {
		return errors.New("amount must be a positive multiple of weekly installment")
	}

	// Get all installments for the loan
	installments, err := uc.loanRepo.GetInstallments(ctx, loanID)
	if err != nil {
		return err
	}

	// Filter unpaid installments that are due (due_date <= today)
	today := time.Now().Truncate(24 * time.Hour)
	var dueUnpaid []models.Installment
	var totalDue float64
	for _, inst := range installments {
		if !inst.Paid && !inst.DueDate.After(today) {
			dueUnpaid = append(dueUnpaid, inst)
			totalDue += inst.Amount

		}
	}

	fmt.Println(dueUnpaid)

	if len(dueUnpaid) > 0 {
		if amount != totalDue {
			return errors.New("payment amount must cover all overdue installments")
		}
	} else {
		// No installments are due – cannot pay ahead
		return errors.New("no installments are due for payment")
	}

	// Prepare IDs of all overdue installments to mark as paid
	var instIDs []int
	for _, inst := range dueUnpaid {
		instIDs = append(instIDs, inst.ID)
	}

	// Create payment record and mark installments as paid
	payment := &models.Payment{
		LoanID:         loanID,
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	}
	err = uc.paymentRepo.Create(ctx, payment, instIDs)
	if err != nil {
		return err
	}

	// Store idempotency key in Redis
	if idempotencyKey != "" {
		_ = uc.idempStore.Store(ctx, idempotencyKey, payment.ID, 24*time.Hour)
	}
	return nil
}
