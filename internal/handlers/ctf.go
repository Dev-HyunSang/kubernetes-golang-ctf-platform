package handlers

import (
	"context"

	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/db"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	flagSubmissions = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ctf_flag_submissions_total",
		Help: "Total number of flag submissions",
	}, []string{"user_id", "status"}) // status: "correct", "incorrect"
)

func SubmitFlag(c *fiber.Ctx) error {
	var input dto.SubmitFlagRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Get user info from JWT
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	// Validate Flag (Dummy logic for now)
	isCorrect := false
	if input.Flag == "CTF{super_secret_flag}" {
		isCorrect = true
	}

	// Save to DB
	_, err := db.Client.Submission.
		Create().
		SetUserID(userID).
		SetFlag(input.Flag).
		SetIsCorrect(isCorrect).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save submission"})
	}

	// Prometheus metric
	status := "incorrect"
	if isCorrect {
		status = "correct"
	}
	flagSubmissions.WithLabelValues(string(rune(userID)), status).Inc()

	if isCorrect {
		return c.JSON(fiber.Map{"message": "Correct flag!"})
	}
	return c.JSON(fiber.Map{"message": "Incorrect flag."})
}
