package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type HospitalHandler struct {
	hospitalUsecase usecase.HospitalUsecase
	config          *config.Config
}

func NewHospitalHandler(hospitalUsecase usecase.HospitalUsecase, cfg *config.Config) *HospitalHandler {
	return &HospitalHandler{
		hospitalUsecase: hospitalUsecase,
		config:          cfg,
	}
}

// GetHospitalMe godoc
// @Summary     Hastane bilgilerini getirir
// @Description Mevcut hastane bilgilerini döner
// @Tags        Hospital
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} dto.HospitalResponse
// @Failure     401 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /api/hospital/me [get]
func (h *HospitalHandler) GetHospitalMe(c *fiber.Ctx) error {
	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	hospital, err := h.hospitalUsecase.GetHospitalByID(user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(hospital)
}

// UpdateHospitalMe godoc
// @Summary     Hastane bilgilerini günceller
// @Description Mevcut hastane bilgilerini günceller
// @Tags        Hospital
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       update body dto.UpdateHospitalRequest true "Update hospital info"
// @Success     200 {object} dto.HospitalResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /api/hospital/me [put]
func (h *HospitalHandler) UpdateHospitalMe(c *fiber.Ctx) error {
	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	req := new(dto.UpdateHospitalRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	hospital, err := h.hospitalUsecase.UpdateHospital(user.HospitalID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(hospital)
}
