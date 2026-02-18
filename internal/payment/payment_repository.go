package payment

import (
	"context"

	"github.com/evrintobing17/loan-billing-system/models"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment, installmentIDs []int) error
	GetByIdempotencyKey(ctx context.Context, key string) (*models.Payment, error)
}
