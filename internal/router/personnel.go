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

func PersonnelRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	personnelRepo := repository.NewPersonnelRepository(db)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, cfg)

	// Location modülü bağımlılıkları
	locationRepo := repository.NewLocationRepository()
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	locationHandler := handler.NewLocationHandler(locationUsecase)

	api := app.Group("/api")

	// Location endpointleri
	api.Get("/cities", locationHandler.ListCities)
	api.Get("/districts", locationHandler.ListDistrictsByCity)

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)
}
