package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/handler"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/middleware"
)

func LocationRoutes(app *fiber.App) {

	locationRepo := repository.NewLocationRepository()
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	locationHandler := handler.NewLocationHandler(locationUsecase)

	api := app.Group("/api")

	locationGroup := api.Group("/location")

	locationGroup.Get("/cities", middleware.GeneralRateLimiter(), locationHandler.ListCities)
	locationGroup.Get("/districts", middleware.GeneralRateLimiter(), locationHandler.ListDistrictsByCity)
}
