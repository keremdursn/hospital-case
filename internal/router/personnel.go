package router

import (
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

func PersonnelRoutes(deps RouterDeps) {
	// db := database.GetDB()
	personnelRepo := repository.NewPersonnelRepository(deps.DB.SQL)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo, deps.DB.Redis)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, deps.Config)

	api := deps.App.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)
}
