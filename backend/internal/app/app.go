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

// New constructs the Fiber application, registers middleware, and wires routes.
func New(_ config.Config) *fiber.App {
	application := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandler,
	})

	application.Use(recover.New())
	application.Use(logger.New())

	health.New().Register(application)

	calculator := calculatorController.New(calculatorService.New())
	calculator.Register(application)

	return application
}
