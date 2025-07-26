package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/usecase"
)

type LocationHandler struct {
	locationUsecase usecase.LocationUsecase
}

func NewLocationHandler(locationUsecase usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{
		locationUsecase: locationUsecase,
	}
}

// ListCities godoc
// @Summary     Tüm şehirleri listeler
// @Description Tüm şehirleri döner
// @Tags        Location
// @Accept      json
// @Produce     json
// @Success     200 {array} dto.CityLookup
// @Router      /api/location/cities [get]
func (h *LocationHandler) ListCities(c *fiber.Ctx) error {
	cities, err := h.locationUsecase.ListAllCities()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(cities)
}

// ListDistrictsByCity godoc
// @Summary     Şehre göre ilçeleri listeler
// @Description Belirli bir şehre ait ilçeleri döner
// @Tags        Location
// @Accept      json
// @Produce     json
// @Param       city_id query int true "City ID"
// @Success     200 {array} dto.DistrictLookup
// @Failure     400 {object} map[string]string
// @Router      /api/location/districts [get]
func (h *LocationHandler) ListDistrictsByCity(c *fiber.Ctx) error {
	cityID, err := strconv.ParseUint(c.Query("city_id", "0"), 10, 64)
	if err != nil || cityID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid city_id"})
	}

	districts, err := h.locationUsecase.ListDistrictsByCity(uint(cityID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(districts)
}
