package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

func SubUserRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	subuserRepo := repository.NewSubUserRepository(db)
	subuserUsecase := usecase.NewSubUserUsecase(subuserRepo)
	subuserHandler := handler.NewSubUserHandler(subuserUsecase, cfg)

	api := app.Group("/api")

	subuserGroup := api.Group("/subuser")

	subuserGroup.Post("/", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), subuserHandler.CreateSubUser)
	subuserGroup.Get("/users", middleware.GeneralRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), subuserHandler.ListUsers)
	subuserGroup.Put("/:id", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), subuserHandler.UpdateSubUser)
	subuserGroup.Delete("/:id", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), subuserHandler.DeleteSubUser)
}
