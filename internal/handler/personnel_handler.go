package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/usecase"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type PersonnelHandler struct {
	usecase usecase.PersonnelUsecase
	config  *config.Config
}

func NewPersonnelHandler(usecase usecase.PersonnelUsecase, cfg *config.Config) *PersonnelHandler {
	return &PersonnelHandler{
		usecase: usecase,
		config:  cfg,
	}
}

func (h *PersonnelHandler) ListAllJobGroups(c *fiber.Ctx) error {
	resp, err := h.usecase.ListAllJobGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *PersonnelHandler) ListTitleByJobGroup(c *fiber.Ctx) error {

	jobGroupID, err := strconv.ParseUint(c.Query("job_group_id", "0"), 10, 64)
	if err != nil || jobGroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job_group_id"})
	}

	resp, err := h.usecase.ListTitleByJobGroup(uint(jobGroupID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *PersonnelHandler) AddStaff(c *fiber.Ctx) error {
	var req dto.AddStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.AddStaff(&req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *PersonnelHandler) UpdateStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid staff id"})
	}

	req := new(dto.UpdateStaffRequest)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.UpdateStaff(uint(id), req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *PersonnelHandler) DeleteStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid staff id"})
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.DeleteStaff(uint(id), user.HospitalID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Staff deleted"})
}

func (h *PersonnelHandler) ListStaff(c *fiber.Ctx) error {

	// 1. Query parametrelerinden filtreleri ve sayfa/size al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	firstName := c.Query("first_name", "")
	lastName := c.Query("last_name", "")
	tc := c.Query("tc", "")

	// job_group_id ve title_id parse edilir
	var jobGroupID *uint
	if v := c.Query("job_group_id", ""); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			jid := uint(id)
			jobGroupID = &jid
		}
	}
	var titleID *uint
	if v := c.Query("title_id", ""); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			tid := uint(id)
			titleID = &tid
		}
	}

	user := utils.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.ListStaff(user.HospitalID, dto.StaffListFilter{
		FirstName:  firstName,
		LastName:   lastName,
		TC:         tc,
		JobGroupID: jobGroupID,
		TitleID:    titleID,
	}, page, size)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}
