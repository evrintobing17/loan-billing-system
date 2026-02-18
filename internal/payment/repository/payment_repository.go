package repository

import (
	"context"
	"database/sql"

	"github.com/evrintobing17/loan-billing-system/internal/payment"
	"github.com/evrintobing17/loan-billing-system/models"
	"github.com/lib/pq"
)

type paymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(DB *sql.DB) payment.PaymentRepository {
	return &paymentRepository{
		DB: DB,
	}
}

// Create implements [payment.PaymentRepository].
func (p *paymentRepository) Create(ctx context.Context, payment *models.Payment, installmentIDs []int) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO payments (loan_id, amount, idempotency_key) 
              VALUES ($1, $2, $3) RETURNING id, payment_date`
	err = tx.QueryRowContext(ctx, query, payment.LoanID, payment.Amount, payment.IdempotencyKey).
		Scan(&payment.ID, &payment.PaymentDate)
	if err != nil {
		return err
	}

	// Mark installments as paid
	if len(installmentIDs) > 0 {
		// Update installments set paid=true, payment_id = ?
		_, err = tx.ExecContext(ctx,
			`UPDATE installments SET paid = true WHERE id = ANY($1)`, pq.Array(installmentIDs))
		if err != nil {
			return err
		}

		// Link payment to installments
		for _, instID := range installmentIDs {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO payment_installments (payment_id, installment_id) VALUES ($1, $2)`,
				payment.ID, instID)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// GetByIdempotencyKey implements [payment.PaymentRepository].
func (p *paymentRepository) GetByIdempotencyKey(ctx context.Context, key string) (*models.Payment, error) {
	panic("unimplemented")
}
