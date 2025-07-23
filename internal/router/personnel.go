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

func JobGroupRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	personnelRepo := repository.NewPersonnelRepository(db)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, cfg)

	api := app.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.DeleteStaff)
}
