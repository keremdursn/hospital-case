package router

import (
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

func HospitalRoutes(deps RouterDeps) {
	hRepo := repository.NewHospitalRepository(deps.DB.SQL)
	hUsecase := usecase.NewHospitalUsecase(hRepo)
	hHandler := handler.NewHospitalHandler(hUsecase, deps.Config)

	api := deps.App.Group("/api")

	hGroup := api.Group("/hospital")

	hGroup.Get("/me", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), hHandler.GetHospitalMe)
	hGroup.Put("/me", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), hHandler.UpdateHospitalMe)
}
