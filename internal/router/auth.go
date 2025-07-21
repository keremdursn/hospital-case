package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

func AuthRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase, cfg)

	api := app.Group("/api")

	authGroup := app.Group("/auth")

	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/forgot-password", authHandler.ForgotPassword)
	authGroup.Post("/reset-password", authHandler.ResetPassword)

	api.Get("/protected", utils.AuthRequired(cfg), func(c *fiber.Ctx) error {
		user := utils.GetUserInfo(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.JSON(fiber.Map{
			"authority_id": user.AuthorityID,
			"hospital_id":  user.HospitalID,
			"role":         user.Role,
		})
	})

	// Sub-user management (only for 'yetkili')
	// api.Post("/users", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), authHandler.CreateSubUser)

}
