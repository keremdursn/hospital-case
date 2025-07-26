package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
)

func AuthRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase, cfg)

	api := app.Group("/api")

	authGroup := api.Group("/auth")

	authGroup.Post("/register", middleware.AuthRateLimiter(), authHandler.Register)
	authGroup.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	authGroup.Post("/forgot-password", middleware.AuthRateLimiter(), authHandler.ForgotPassword)
	authGroup.Post("/reset-password", middleware.AuthRateLimiter(), authHandler.ResetPassword)
	authGroup.Post("/refresh-token", middleware.AuthRateLimiter(), authHandler.RefreshToken)
}
