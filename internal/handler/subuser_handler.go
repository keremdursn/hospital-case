package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type SubUserHandler struct {
	usecase usecase.SubUserUsecase
}

func NewSubUserHandler(usecase usecase.SubUserUsecase) *SubUserHandler {
	return &SubUserHandler{usecase: usecase}
}

func (h *SubUserHandler) CreateSubUser(c *fiber.Ctx) error {
	req := new(usecase.CreateSubUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.CreateSubUser(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
