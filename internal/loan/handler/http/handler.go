package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/evrintobing17/loan-billing-system/internal/loan"
	"github.com/evrintobing17/loan-billing-system/models"
	"github.com/gin-gonic/gin"
)

type LoanHandler struct {
	loanUC loan.LoanUsecase
}

func NewLoanHandler(uc loan.LoanUsecase) *LoanHandler {
	return &LoanHandler{loanUC: uc}
}

// CreateLoan godoc
// @Summary Create a new loan
// @Description Create a loan with given terms. Generates weekly installments.
// @Tags loans
// @Accept json
// @Produce json
// @Param request body models.CreateLoanRequest true "Loan details"
// @Success 201 {object} models.Loan
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loans [post]
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req models.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}

	loan, err := h.loanUC.CreateLoan(c.Request.Context(), req.Principal, req.InterestRate, req.TermWeeks, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loan)
}

// GetOutstanding godoc
// @Summary Get outstanding amount for a loan
// @Tags loans
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]float64
// @Failure 404 {object} map[string]string
// @Router /loans/{id}/outstanding [get]
func (h *LoanHandler) GetOutstanding(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}
	outstanding, err := h.loanUC.GetOutstanding(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"outstanding": outstanding})
}

// IsDelinquent godoc
// @Summary Check if a borrower is delinquent
// @Tags loans
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]bool
// @Failure 404 {object} map[string]string
// @Router /loans/{id}/delinquent [get]
func (h *LoanHandler) IsDelinquent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan id"})
		return
	}
	delinquent, err := h.loanUC.IsDelinquent(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"delinquent": delinquent})
}
