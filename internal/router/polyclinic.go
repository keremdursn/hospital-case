package router

import (
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

func PolyclinicRoutes(deps RouterDeps) {
	// db := database.GetDB()
	polyclinicRepo := repository.NewPolyclinicRepository(deps.DB.SQL)
	polyclinicUsecase := usecase.NewPolyclinicUsecase(polyclinicRepo)
	polyclinicHandler := handler.NewPolyclinicHandler(polyclinicUsecase, deps.Config)

	api := deps.App.Group("/api")

	polyclinicGroup := api.Group("/polyclinic")

	polyclinicGroup.Get("/", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListAllPolyclinics)
	polyclinicGroup.Post("/hospital-polyclinics", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), polyclinicHandler.AddHospitalPolyclinic)
	polyclinicGroup.Get("/hospital-polyclinics", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListHospitalPolyclinic)
	polyclinicGroup.Delete("/hospital-polyclinics/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), polyclinicHandler.RemoveHospitalPolyclinic)
}
