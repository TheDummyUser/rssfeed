package handlers

import (
	"time"

	"github.com/TheDummyUser/goRss/config"
	"github.com/TheDummyUser/goRss/model"
	"github.com/TheDummyUser/goRss/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx, db *gorm.DB) error {
	var input model.LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user model.User

	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "invalid username or password",
			"status_code": fiber.StatusUnauthorized,
		})
	}

	if !utils.ComparePassword(user.Password, input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "invalid username or password",
			"status_code": fiber.StatusUnauthorized,
		})
	}

	// Generate Access Token (short-lived)
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["username"] = user.Username
	accessClaims["user_id"] = user.ID
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix() // Expires in 15 min

	aToken, err := accessToken.SignedString([]byte(config.Config("TOKEN")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Generate Refresh Token (long-lived)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = user.ID
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Expires in 7 days

	rToken, err := refreshToken.SignedString([]byte(config.Config("TOKEN")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Optional: Save refresh token in DB
	user.RefreshToken = &rToken
	db.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "User logged in successfully",
		"status_code": fiber.StatusOK,
		"details": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"tokens": fiber.Map{
				"access_token":  aToken,
				"refresh_token": rToken,
			},
		},
	})
}

func Signup(c *fiber.Ctx, db *gorm.DB) error {

	var input model.SignupRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var existingUser model.User
	if err := db.Where("email = ? OR username = ?", input.Email, input.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email or username already exists",
			"status_code": fiber.StatusConflict,
		})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user", "status_code": fiber.StatusInternalServerError})
	}

	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user", "status_code": fiber.StatusInternalServerError})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "User created successfully",
		"statusCode": fiber.StatusOK,
		"details": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

func Logout(c *fiber.Ctx, db *gorm.DB) error {
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

	db.Model(&model.User{}).Where("id = ?", userID).Update("refresh_token", "")

	return c.JSON(fiber.Map{"message": "Logged out successfully", "status_code": fiber.StatusOK})
}
