package main

import (
	"log"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/db"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/handlers"
	"github.com/dev-hyunsang/kubernetes-golang-ctf-platform/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize Database and Redis
	db.InitDatabase()
	db.InitRedis()
	defer db.Client.Close()

	app := fiber.New()
	app.Use(logger.New())

	// Prometheus Middleware
	prometheus := fiberprometheus.New("ctf_platform")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// API Routes
	api := app.Group("/api/v1")

	// Public Routes
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)

	// Protected Routes
	protected := api.Group("/", middleware.Protected())

	// CTF Routes (Users, Admins, SuperAdmins)
	protected.Post("/submit", middleware.RequireRole("user", "admin", "superadmin"), handlers.SubmitFlag)

	// Admin Route Example
	protected.Get("/admin/stats", middleware.RequireRole("admin", "superadmin"), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin dashboard stats"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
