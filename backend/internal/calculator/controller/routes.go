package controller

import "github.com/gofiber/fiber/v3"

// Register mounts calculator routes on the given router.
func (c *Controller) Register(router fiber.Router) {
	api := router.Group("/api/v1")
	api.Post("/evaluate", c.Evaluate)
}
