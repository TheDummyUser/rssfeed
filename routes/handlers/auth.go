package handlers

import (
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

	if err := db.Where("username", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "invalid email or password",
			"status_code": fiber.StatusUnauthorized,
		})
	}

	if !utils.ComparePassword(user.Password, input.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "invalid username or password",
			"status_code": fiber.StatusUnauthorized,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID

	t, err := token.SignedString([]byte(config.Config("TOKEN")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user login sucesssfully", "status_code": fiber.StatusOK,
		"details": fiber.Map{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"token":      t,
		}})
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
