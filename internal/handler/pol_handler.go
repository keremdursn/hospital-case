package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type PolyclinicHandler struct {
	polyclinicUsecase usecase.PolyclinicUsecase
	config            *config.Config
}

func NewPolyclinicHandler(polyclinicUsecase usecase.PolyclinicUsecase, cfg *config.Config) *PolyclinicHandler {
	return &PolyclinicHandler{
		polyclinicUsecase: polyclinicUsecase,
		config:            cfg,
	}
}

// ListAllPolyclinics godoc
// @Summary     Tüm poliklinikleri listeler
// @Description Tüm poliklinikleri döner
// @Tags        Polyclinic
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} dto.PolyclinicLookup
// @Router      /api/polyclinic/ [get]
func (h *PolyclinicHandler) ListAllPolyclinics(c *fiber.Ctx) error {
	polyclinics, err := h.polyclinicUsecase.ListAllPolyclinics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(polyclinics)
}

// AddHospitalPolyclinic godoc
// @Summary     Hastaneye poliklinik ekler
// @Description Hastaneye yeni poliklinik ekler
// @Tags        Polyclinic
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       hospital_polyclinic body dto.AddHospitalPolyclinicRequest true "Add hospital polyclinic info"
// @Success     201 {object} dto.HospitalPolyclinicResponse
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/add [post]
func (h *PolyclinicHandler) AddHospitalPolyclinic(c *fiber.Ctx) error {
	req := new(dto.AddHospitalPolyclinicRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.polyclinicUsecase.AddPolyclinicToHospital(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListHospitalPolyclinic godoc
// @Summary     Hastanenin polikliniklerini listeler (sayfalı)
// @Description Lists hospital's polyclinics with pagination
// @Tags        Polyclinic
// @Accept      json
// @Produce     json
// @Param       page query int false "Page number"
// @Param       size query int false "Page size"
// @Security    BearerAuth
// @Success     200 {object} dto.HospitalPolyclinicListResponse
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/list [get]
func (h *PolyclinicHandler) ListHospitalPolyclinic(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.polyclinicUsecase.ListHospitalPolyclinic(user.HospitalID, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// RemoveHospitalPolyclinic godoc
// @Summary     Hastane polikliniğini siler
// @Description Belirtilen hastane polikliniğini siler
// @Tags        Polyclinic
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Hospital Polyclinic ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/hospital-polyclinics/{id} [delete]
func (h *PolyclinicHandler) RemoveHospitalPolyclinic(c *fiber.Ctx) error {
	hospitalPolyclinicID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = h.polyclinicUsecase.RemoveHospitalPolyclinic(uint(hospitalPolyclinicID), user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Hospital polyclinic removed successfully"})
}
