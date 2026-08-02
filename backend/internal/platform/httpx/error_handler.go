package httpx

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"pascal-tech-dev/calczzle-backend/internal/calculator/expression"
)

// ErrorHandler maps domain and transport errors to consistent JSON responses.
func ErrorHandler(ctx fiber.Ctx, err error) error {
	status, code, message := mapError(err)
	return ctx.Status(status).JSON(ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func mapError(err error) (status int, code, message string) {
	switch {
	case errors.Is(err, expression.ErrDivisionByZero):
		return fiber.StatusUnprocessableEntity, "DIVISION_BY_ZERO", "Division by zero is not allowed."
	case errors.Is(err, expression.ErrInvalidSquareRoot):
		return fiber.StatusUnprocessableEntity, "INVALID_SQUARE_ROOT", "Square root of a negative number is not allowed."
	case errors.Is(err, expression.ErrUnsupportedFunction):
		return fiber.StatusBadRequest, "UNSUPPORTED_FUNCTION", "The function is not supported."
	case errors.Is(err, expression.ErrInvalidExpression):
		return fiber.StatusBadRequest, "INVALID_EXPRESSION", "The expression is invalid."
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		switch fiberErr.Code {
		case fiber.StatusNotFound:
			return fiber.StatusNotFound, "NOT_FOUND", "The requested resource was not found."
		case fiber.StatusRequestEntityTooLarge:
			return fiber.StatusRequestEntityTooLarge, "REQUEST_TOO_LARGE", "The request body is too large."
		default:
			return fiberErr.Code, "INTERNAL_SERVER_ERROR", "An unexpected error occurred."
		}
	}

	// Tokenizer/parser still return plain errors for syntax problems.
	return fiber.StatusBadRequest, "INVALID_EXPRESSION", "The expression is invalid."
}
