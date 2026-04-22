package handlers

import (
	"context"
	"time"

	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/ent/user"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/db"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/dto"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	u, err := db.Client.User.
		Create().
		SetEmail(input.Email).
		SetPassword(string(hashedPassword)).
		SetNickname(input.Nickname).
		SetAffiliation(input.Affiliation).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    fiber.Map{"id": u.ID, "email": u.Email, "nickname": u.Nickname},
	})
}

func Login(c *fiber.Ctx) error {
	var input dto.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	u, err := db.Client.User.Query().Where(user.Email(input.Email)).Only(context.Background())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := utils.GenerateToken(u.ID, u.Email, string(u.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Optional: Store token session in Redis (e.g., for blacklisting)
	err = db.RedisClient.Set(context.Background(), "session:"+token, u.ID, 72*time.Hour).Err()
	if err != nil {
		// Log error, but proceed
	}

	return c.JSON(fiber.Map{"token": token})
}
