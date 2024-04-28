package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	middlewares "github.com/Cannonskr/assessment-tax/src/api/middlewares"
	service "github.com/Cannonskr/assessment-tax/src/api/services"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	// user
	e.POST("/tax/calculations", service.TaxCalculationHandler)

	// admin
	adminAuth := middlewares.AdminAuthMiddleware()
	// Authentication (admin)
	adminGroup := e.Group("/admin", adminAuth)

	// Endpoint update "personal"
	adminGroup.POST("/deductions/personal", service.UpdatePersonalDeductionHandler)

	// Endpoint update "k-receipt"
	adminGroup.POST("/deductions/k-receipt", service.UpdateKReceiptDeductionHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Printf("Shutting down the server: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Error during graceful shutdown: %v", err)
	}
}
