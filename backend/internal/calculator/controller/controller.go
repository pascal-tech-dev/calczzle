package controller

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"pascal-tech-dev/calczzle-backend/internal/calculator/service"
)

// Controller serves calculator endpoints.
type Controller struct {
	evaluator service.Evaluator
}

// New constructs a calculator Controller.
func New(evaluator service.Evaluator) *Controller {
	return &Controller{evaluator: evaluator}
}

// Evaluate handles POST /api/v1/evaluate.
func (c *Controller) Evaluate(ctx fiber.Ctx) error {
	var req EvaluateRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorBody{
			Error: errorDetail{
				Code:    "INVALID_REQUEST",
				Message: "The request body is invalid.",
			},
		})
	}

	expression := strings.TrimSpace(req.Expression)
	if expression == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(errorBody{
			Error: errorDetail{
				Code:    "EMPTY_EXPRESSION",
				Message: "The expression must not be empty.",
			},
		})
	}

	result, err := c.evaluator.Evaluate(expression)
	if err != nil {
		return err
	}

	return ctx.JSON(EvaluateResponse{Result: result})
}
