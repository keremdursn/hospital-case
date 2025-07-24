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

func HospitalRoutes(app *fiber.App, cfg *config.Config) {
	db := database.GetDB()
	hRepo := repository.NewHospitalRepository(db)
	hUsecase := usecase.NewHospitalUsecase(hRepo)
	hHandler := handler.NewHospitalHandler(hUsecase, cfg)

	api := app.Group("/api")

	hGroup := api.Group("/hospital")

	hGroup.Get("/hospital/me", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), hHandler.GetHospitalMe)
	hGroup.Put("/hospital/me", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), hHandler.UpdateHospitalMe)
}
