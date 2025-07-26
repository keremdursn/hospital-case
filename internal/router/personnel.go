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

func PersonnelRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	personnelRepo := repository.NewPersonnelRepository(db)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, cfg)

	api := app.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", middleware.GeneralRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", middleware.GeneralRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", middleware.GeneralRateLimiter(), utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)
}
