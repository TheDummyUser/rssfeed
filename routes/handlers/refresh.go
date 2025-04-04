package handlers

import (
	"time"

	"github.com/TheDummyUser/goRss/config"
	"github.com/TheDummyUser/goRss/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RefreshToken(c *fiber.Ctx, db *gorm.DB) error {
	var req struct {
		RefreshToken string `json:"r_token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Find user by refresh token
	var user model.User
	if err := db.Where("refresh_token = ?", req.RefreshToken).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Generate new Access Token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // New 15-min token

	aToken, err := accessToken.SignedString([]byte(config.Config("TOKEN")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"a_token": aToken})
}
