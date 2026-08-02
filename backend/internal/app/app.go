package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	calculatorController "pascal-tech-dev/calczzle-backend/internal/calculator/controller"
	calculatorService "pascal-tech-dev/calczzle-backend/internal/calculator/service"
	"pascal-tech-dev/calczzle-backend/internal/config"
	"pascal-tech-dev/calczzle-backend/internal/health"
	"pascal-tech-dev/calczzle-backend/internal/platform/httpx"
)

// MaxRequestBodyBytes caps JSON request bodies for evaluate (and other) endpoints.
// A few kilobytes is more than enough for calculator expressions.
const MaxRequestBodyBytes = 4 * 1024

// New constructs the Fiber application, registers middleware, and wires routes.
func New(_ config.Config) *fiber.App {
	application := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandler,
		BodyLimit:    MaxRequestBodyBytes,
	})

	application.Use(recover.New())
	application.Use(logger.New())

	health.New().Register(application)

	calculator := calculatorController.New(calculatorService.New())
	calculator.Register(application)

	return application
}
