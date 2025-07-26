package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type SubUserHandler struct {
	subUserUsecase usecase.SubUserUsecase
	config         *config.Config
}

func NewSubUserHandler(subUserUsecase usecase.SubUserUsecase, cfg *config.Config) *SubUserHandler {
	return &SubUserHandler{
		subUserUsecase: subUserUsecase,
		config:         cfg,
	}
}

// CreateSubUser godoc
// @Summary     Alt kullanıcı oluşturur
// @Description Yeni alt kullanıcı kaydı oluşturur (yetkili/çalışan)
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       subuser body dto.CreateSubUserRequest true "Create subuser info"
// @Success     201 {object} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser [post]
func (h *SubUserHandler) CreateSubUser(c *fiber.Ctx) error {
	req := new(dto.CreateSubUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.subUserUsecase.CreateSubUser(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListUsers godoc
// @Summary     Alt kullanıcıları listeler
// @Description Hastaneye ait tüm kullanıcıları listeler
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser [get]
func (h *SubUserHandler) ListUsers(c *fiber.Ctx) error {
	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.subUserUsecase.ListUsers(user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateSubUser godoc
// @Summary     Alt kullanıcı bilgilerini günceller
// @Description Mevcut alt kullanıcı bilgilerini günceller
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "SubUser ID"
// @Param       subuser body dto.UpdateSubUserRequest true "Update subuser info"
// @Success     200 {object} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser/{id} [put]
func (h *SubUserHandler) UpdateSubUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	req := new(dto.UpdateSubUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.subUserUsecase.UpdateSubUser(uint(id), req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteSubUser godoc
// @Summary     Alt kullanıcı siler
// @Description Belirtilen alt kullanıcıyı siler
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "SubUser ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/subuser/{id} [delete]
func (h *SubUserHandler) DeleteSubUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.subUserUsecase.DeleteSubUser(uint(id), user.HospitalID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "SubUser deleted successfully"})
}
