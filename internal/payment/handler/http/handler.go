package http

import (
	"net/http"
	"strconv"

	"github.com/evrintobing17/loan-billing-system/internal/payment"
	"github.com/evrintobing17/loan-billing-system/models"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUC payment.PaymentUsecase
}

func NewPaymentHandler(uc payment.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUC: uc}
}

// MakePayment godoc
// @Summary Make a payment against a loan
// @Description Process a payment. Idempotency-Key header prevents duplicates.
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param request body models.PaymentRequest true "Payment amount"
// @Param Idempotency-Key header string true "Unique idempotency key"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loans/{id}/payments [post]
func (h *PaymentHandler) MakePayment(c *gin.Context) {
	var req models.PaymentRequest
	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
		return
	}

	err = h.paymentUC.MakePayment(c.Request.Context(), loanID, req.Amount, idempotencyKey)
	if err != nil {
		if err.Error() == "amount must be a positive multiple of weekly installment" ||
			err.Error() == "amount exceeds total outstanding" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment processed"})
}
