package api

import (
	"github.com/TheDummyUser/goRss/routes"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *fiber.App {
	app := fiber.New()

	routes.SetupRoutes(app, db)

	return app
}
