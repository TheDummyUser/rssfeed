package routes

import (
	"github.com/TheDummyUser/goRss/middleware"
	"github.com/TheDummyUser/goRss/routes/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	// login and signup
	api.Post("/login", func(c *fiber.Ctx) error {
		return handlers.Login(c, db)
	})
	api.Post("/register", func(c *fiber.Ctx) error {
		return handlers.Signup(c, db)
	})
	api.Get("/logout", func(c *fiber.Ctx) error {
		return handlers.Logout(c, db)
	})

	api.Post("/refresh", func(c *fiber.Ctx) error {
		return handlers.RefreshToken(c, db)
	})
	protected := api.Group("", middleware.Protected())

	// rss add and rss view
	protected.Post("/rss_add", func(c *fiber.Ctx) error {
		return handlers.RssAdd(c, db)
	})

	protected.Get("/fetchrss", func(c *fiber.Ctx) error {
		return handlers.FetchRssFeeds(c, db)
	})

	protected.Post("/summarize", func(c *fiber.Ctx) error {
		return handlers.SummarizeAi(c, db)
	})

}
