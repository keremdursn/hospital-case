package router

import (
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
)

func AuthRoutes(deps RouterDeps) {
	// db := database.GetDB()
	authRepo := repository.NewAuthRepository(deps.DB.SQL)
	authUsecase := usecase.NewAuthUsecase(authRepo, deps.DB.Redis)
	authHandler := handler.NewAuthHandler(authUsecase, deps.Config)

	api := deps.App.Group("/api")

	authGroup := api.Group("/auth")

	authGroup.Post("/register", middleware.AuthRateLimiter(), authHandler.Register)
	authGroup.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	authGroup.Post("/forgot-password", middleware.AuthRateLimiter(), authHandler.ForgotPassword)
	authGroup.Post("/reset-password", middleware.AuthRateLimiter(), authHandler.ResetPassword)
	authGroup.Post("/refresh-token", middleware.AuthRateLimiter(), authHandler.RefreshToken)
}
