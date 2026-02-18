package main

import (
	"log"

	"github.com/evrintobing17/loan-billing-system/config"
	loanHttp "github.com/evrintobing17/loan-billing-system/internal/loan/handler/http"
	loanRepo "github.com/evrintobing17/loan-billing-system/internal/loan/repository"
	loanUsecase "github.com/evrintobing17/loan-billing-system/internal/loan/usecase"
	paymentHttp "github.com/evrintobing17/loan-billing-system/internal/payment/handler/http"
	paymentRepo "github.com/evrintobing17/loan-billing-system/internal/payment/repository"
	paymentUsecase "github.com/evrintobing17/loan-billing-system/internal/payment/usecase"
	"github.com/evrintobing17/loan-billing-system/pkg/idempotency"
	redisClient "github.com/evrintobing17/loan-billing-system/pkg/redis"
	"github.com/joho/godotenv"

	"github.com/evrintobing17/loan-billing-system/pkg/postgres"

	_ "github.com/evrintobing17/loan-billing-system/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Loan Billing API
// @version 1.0
// @description Loan billing system with weekly installments.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	// PostgreSQL connection
	db, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Redis connection
	rdb, err := redisClient.NewClientWithRetry(cfg.RedisAddr, 5)
	if err != nil {
		log.Fatal("Failed to connect to Redis after retries:", err)
	}

	// Initialize repositories
	lRepo := loanRepo.NewLoanRepository(db)
	pRepo := paymentRepo.NewPaymentRepository(db)

	// Idempotency store
	idempStore := idempotency.NewRedisStore(rdb)

	// Use cases
	loanUC := loanUsecase.NewLoanUseCase(lRepo)
	paymentUC := paymentUsecase.NewPaymentUseCase(pRepo, lRepo, idempStore)

	// Handlers
	loanHandler := loanHttp.NewLoanHandler(loanUC)
	paymentHandler := paymentHttp.NewPaymentHandler(paymentUC)

	// Gin engine
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/loans", loanHandler.CreateLoan)
		v1.GET("/loans/:id/outstanding", loanHandler.GetOutstanding)
		v1.GET("/loans/:id/delinquent", loanHandler.IsDelinquent)
		v1.POST("/loans/:id/payments", paymentHandler.MakePayment)
	}

	r.Run(":" + cfg.Port)
}
