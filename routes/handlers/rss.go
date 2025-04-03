package handlers

import (
	"github.com/TheDummyUser/goRss/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

func RssAdd(c *fiber.Ctx, db *gorm.DB) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user token"})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user_id missing or invalid"})
	}
	userID := uint(userIDFloat)
	var inputs model.RssRequest

	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var existingRssUrl model.RssData

	if err := db.Where("name = ? or url = ?", inputs.Name, inputs.Url).First(&existingRssUrl).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status_code": fiber.StatusConflict,
			"error":       "url or name already exists",
		})
	}

	rssurl := model.RssData{
		UserID: userID,
		Name:   inputs.Name,
		Url:    inputs.Url,
	}

	if err := db.Create(&rssurl).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add rss url", "status_code": fiber.StatusInternalServerError, "main_error": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "url added sucesssfully", "status_code": fiber.StatusOK,
		"details": fiber.Map{
			"id":   rssurl.ID,
			"name": rssurl.Name,
			"url":  rssurl.Url,
		}})
}

func FetchRssFeeds(c *fiber.Ctx, db *gorm.DB) error {

	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user token"})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user_id missing or invalid"})
	}
	userID := uint(userIDFloat)

	// Retrieve all RSS URLs added by the user
	var rssFeeds []model.RssData
	if err := db.Where("user_id = ?", userID).Find(&rssFeeds).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch RSS URLs"})
	}

	// Initialize RSS parser
	parser := gofeed.NewParser()
	var fetchedFeeds []fiber.Map

	// Fetch and parse each RSS feed
	for _, feedData := range rssFeeds {
		feed, err := parser.ParseURL(feedData.Url)
		if err != nil {
			continue // Skip if parsing fails
		}

		var feedItems []fiber.Map
		for _, item := range feed.Items {
			// Handle nil values to prevent panic
			var imageURL, authorEmail, authorName string
			if item.Image != nil {
				imageURL = item.Image.URL
			}
			if item.Author != nil {
				authorEmail = item.Author.Email
				authorName = item.Author.Name
			}

			feedItems = append(feedItems, fiber.Map{
				"title":       item.Title,
				"link":        item.Link,
				"description": item.Description,
				"published":   item.Published,
				"image":       imageURL,
				"content":     item.Content,
				"links":       item.Links,
				"authors":     item.Authors,
				"authorMail":  authorEmail,
				"authorName":  authorName,
			})
		}

		fetchedFeeds = append(fetchedFeeds, fiber.Map{
			"name":  feedData.Name,
			"url":   feedData.Url,
			"title": feed.Title,
			"items": feedItems,
		})
	}

	// Return the fetched RSS feeds
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"feeds":       fetchedFeeds,
	})
}
