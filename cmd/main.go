package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/config"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/delivery/rest"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/database"
	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	config := config.NewConfig()

	// init db
	db, err := database.InitDB(config)
	if err != nil {
		log.Panicf("Failed connect to database: %v", err)
	}

	// init repository
	consumerRepo := repository.NewConsumerRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	consumerLimitRepo := repository.NewConsumerLimitRepository(db)
	loanRepo := repository.NewLoanRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// init usecase
	consumerUC := usecase.NewConsumerUsecase(consumerRepo, config.Timeout)
	merchantUC := usecase.NewMerchantUsecase(merchantRepo, config.Timeout)
	consumerLimitUC := usecase.NewConsumerLimitUsecase(consumerLimitRepo, config.Timeout)
	loanUC := usecase.NewLoanUsecase(
		loanRepo,
		consumerLimitRepo,
		consumerRepo,
		merchantRepo,
		config.Timeout,
	)
	transactionUC := usecase.NewTransactionUsecase(
		transactionRepo,
		loanRepo,
		consumerLimitRepo,
		consumerRepo,
		config.Timeout,
	)

	// init global middleware
	e.Use(middleware.LoggerMiddleware())
	e.Use(middleware.CORSMiddleware())

	// init handler
	v1 := e.Group("/api/v1")
	rest.NewConsumerHandler(v1, consumerUC)
	rest.NewMerchantHandler(v1, merchantUC)
	rest.NewConsumerLimitHandler(v1, consumerLimitUC)
	rest.NewLoanHandler(v1, loanUC)
	rest.NewTransactionHandler(v1, transactionUC)

	e.Logger.Fatal(e.Start(":" + config.Port))

	// Start server
	go func() {
		if err := e.Start(":" + config.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
