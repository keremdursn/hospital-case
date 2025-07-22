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
	usecase usecase.PolyclinicUsecase
	config  *config.Config
}

func NewPolyclinicHandler(usecase usecase.PolyclinicUsecase, cfg *config.Config) *PolyclinicHandler {
	return &PolyclinicHandler{
		usecase: usecase,
		config:  cfg,
	}
}

func (h *PolyclinicHandler) ListAllPolyclinics(c *fiber.Ctx) error {
	resp, err := h.usecase.ListAllPolyclinics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PolyclinicHandler) AddHospitalPolyclinic(c *fiber.Ctx) error {
	req := new(dto.AddHospitalPolyclinicRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.AddPolyclinicToHospital(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *PolyclinicHandler) ListHospitalPolyclinic(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.ListHospitalPolyclinic(user.HospitalID, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PolyclinicHandler) RemoveHospitalPolyclinic(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = h.usecase.RemoveHospitalPolyclinic(uint(id), user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Polyclinic removed from hospital"})
}
