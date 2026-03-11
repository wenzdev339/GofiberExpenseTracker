package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"gofiber-expense-tracker/internal/handlers"
	"gofiber-expense-tracker/internal/middleware"
	"gofiber-expense-tracker/internal/repositories"
	"gofiber-expense-tracker/internal/services"
)

func Setup(app *fiber.App, db *sql.DB, jwtSecret string, jwtExpireHours int) {
	// Repositories
	transactionRepo := repositories.NewTransactionRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// Services
	transactionSvc := services.NewTransactionService(transactionRepo)
	authSvc := services.NewAuthService(userRepo, jwtSecret, jwtExpireHours)

	// Handlers
	transactionHandler := handlers.NewTransactionHandler(transactionSvc)
	authHandler := handlers.NewAuthHandler(authSvc)

	api := app.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Transaction routes (protected)
	transactions := api.Group("/transactions", middleware.JWTProtected(jwtSecret))
	transactions.Post("/", transactionHandler.Create)
	transactions.Get("/", transactionHandler.List)
	transactions.Get("/summary", transactionHandler.GetSummary)
	transactions.Get("/:id", transactionHandler.GetByID)
	transactions.Put("/:id", transactionHandler.Update)
	transactions.Delete("/:id", transactionHandler.Delete)
}
