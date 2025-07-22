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

func PolyclinicRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	polyclinicRepo := repository.NewPolyclinicRepository(db)
	polyclinicUsecase := usecase.NewPolyclinicUsecase(polyclinicRepo)
	polyclinicHandler := handler.NewPolyclinicHandler(polyclinicUsecase, cfg)

	api := app.Group("/api")

	polyclinicGroup := api.Group("/polyclinic")

	//
	polyclinicGroup.Get("/", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListAllPolyclinics)
	polyclinicGroup.Post("/hospital-polyclinics", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), polyclinicHandler.AddHospitalPolyclinic)
	polyclinicGroup.Get("/hospital-polyclinics", utils.AuthRequired(cfg), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListHospitalPolyclinic)
	polyclinicGroup.Delete("/hospital-polyclinics/:id", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), polyclinicHandler.RemoveHospitalPolyclinic)
}
