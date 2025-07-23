package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/usecase"
)

type LocationHandler struct {
	usecase usecase.LocationUsecase
}

func NewLocationHandler(usecase usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: usecase}
}

// GET /api/cities
func (h *LocationHandler) ListCities(c *fiber.Ctx) error {
	resp, err := h.usecase.ListAllCities()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// GET /api/districts?city_id=1
func (h *LocationHandler) ListDistrictsByCity(c *fiber.Ctx) error {
	cityID, err := strconv.ParseUint(c.Query("city_id", "0"), 10, 64)
	if err != nil || cityID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid city_id"})
	}
	resp, err := h.usecase.ListDistrictsByCity(uint(cityID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
