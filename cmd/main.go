package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/router"

	_ "github.com/keremdursn/hospital-case/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/keremdursn/hospital-case/pkg/metrics"
)

func main() {

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Connect to database
	database.Connect(&cfg)
	database.ConnectRedis(&cfg)

	err = database.DB.AutoMigrate(
		&models.City{},
		&models.District{},
		&models.Hospital{},
		&models.Authority{},
		&models.Polyclinic{},
		&models.HospitalPolyclinic{},
		&models.JobGroup{},
		&models.Title{},
		&models.Staff{},
	)
	if err != nil {
		log.Fatal("cannot migrate database: ", err)
	}

	app := fiber.New()

	app.Use(metrics.PrometheusMiddleware())
	app.Get("/metrics", metrics.PrometheusHandler())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	router.AuthRoutes(app, &cfg)
	router.HospitalRoutes(app, &cfg)
	router.SubUserRoutes(app, &cfg)
	router.PolyclinicRoutes(app, &cfg)
	router.PersonnelRoutes(app, &cfg)
	router.LocationRoutes(app)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
