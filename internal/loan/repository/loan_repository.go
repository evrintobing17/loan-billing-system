package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evrintobing17/loan-billing-system/internal/loan"
	"github.com/evrintobing17/loan-billing-system/models"
	"github.com/lib/pq"
)

type loanRepository struct {
	DB *sql.DB
}

func NewLoanRepository(DB *sql.DB) loan.LoanRepository {
	return &loanRepository{
		DB: DB,
	}
}

// Create implements [loan.LoanRepository].
func (l *loanRepository) Create(ctx context.Context, loan *models.Loan, installments []models.Installment) error {
	tx, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert loan
	query := `INSERT INTO loans (principal, interest_rate, term_weeks, weekly_amount, start_date, is_active) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, query, loan.Principal, loan.InterestRate, loan.TermWeeks,
		loan.WeeklyAmount, loan.StartDate, true).Scan(&loan.ID, &loan.CreatedAt)
	if err != nil {
		return err
	}

	// Insert installments
	for _, inst := range installments {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO installments (loan_id, week_number, due_date, amount) VALUES ($1, $2, $3, $4)`,
			loan.ID, inst.WeekNumber, inst.DueDate, inst.Amount)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (l *loanRepository) GetByID(ctx context.Context, id int) (*models.Loan, error) {
	var loan models.Loan
	query := `SELECT id, principal, interest_rate, term_weeks, weekly_amount, start_date, is_active, created_at 
              FROM loans WHERE id = $1 AND is_active IS TRUE`
	err := l.DB.QueryRowContext(ctx, query, id).Scan(
		&loan.ID,
		&loan.Principal,
		&loan.InterestRate,
		&loan.TermWeeks,
		&loan.WeeklyAmount,
		&loan.StartDate,
		&loan.IsActive,
		&loan.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err // caller can check with errors.Is
		}
		return nil, fmt.Errorf("query loan by id: %w", err)
	}
	return &loan, nil
}

// GetInstallments retrieves all installments for a given loan, ordered by week_number.
func (l *loanRepository) GetInstallments(ctx context.Context, loanID int) ([]models.Installment, error) {
	query := `SELECT id, loan_id, week_number, due_date, amount, paid 
              FROM installments 
              WHERE loan_id = $1 
              ORDER BY week_number`
	rows, err := l.DB.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("query installments: %w", err)
	}
	defer rows.Close()

	var installments []models.Installment
	for rows.Next() {
		var inst models.Installment
		err := rows.Scan(
			&inst.ID,
			&inst.LoanID,
			&inst.WeekNumber,
			&inst.DueDate,
			&inst.Amount,
			&inst.Paid,
		)
		if err != nil {
			return nil, fmt.Errorf("scan installment: %w", err)
		}
		installments = append(installments, inst)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return installments, nil
}

// UpdateInstallmentsPaid marks the given installment IDs as paid.
// This method is typically called within a transaction from the payment repository.
func (l *loanRepository) UpdateInstallmentsPaid(ctx context.Context, installmentIDs []int) error {
	if len(installmentIDs) == 0 {
		return nil // nothing to do
	}

	query := `UPDATE installments SET paid = true WHERE id = ANY($1)`
	_, err := l.DB.ExecContext(ctx, query, pq.Array(installmentIDs))
	if err != nil {
		return fmt.Errorf("update installments paid: %w", err)
	}
	return nil
}
