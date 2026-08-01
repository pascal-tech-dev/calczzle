package health

import "github.com/gofiber/fiber/v3"

// Controller serves infrastructure health endpoints.
type Controller struct{}

// New constructs a health Controller.
func New() *Controller {
	return &Controller{}
}

// Register mounts health routes on the given router.
func (c *Controller) Register(router fiber.Router) {
	router.Get("/health", c.Check)
}

// Check returns a simple readiness response.
func (c *Controller) Check(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}
