package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
)

type RouterDeps struct {
	App    *fiber.App
	DB     *database.Database
	Config *config.Config
}
