package example

import (
	appctx "wallet-api/app/ctx"
	"wallet-api/app/presentation"
	"wallet-api/app/services/example"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Singleton validator instance for better performance
var validate = validator.New()

type exampleHandler struct {
	exampleService example.ExampleService
}

func NewExampleHandler(exampleService example.ExampleService) ExampleHandler {
	return &exampleHandler{
		exampleService: exampleService,
	}
}

func (h *exampleHandler) Create(c *fiber.Ctx) error {
	requestID := appctx.GetRequestID(c.UserContext())
	if requestID == "" {
		if val := c.Locals("request_id"); val != nil {
			if str, ok := val.(string); ok {
				requestID = str
			}
		}
	}

	var req presentation.ExampleRequest
	if err := c.BodyParser(&req); err != nil {
		// Don't expose raw parsing error details
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Invalid request body", requestID))
	}

	if err := validate.Struct(req); err != nil {
		// Return user-friendly validation message without exposing field details
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Validation failed: please check your input", requestID))
	}

	// Pass context.Context to service instead of fiber.Ctx
	ctx := appctx.WithRequestID(c.UserContext(), requestID)
	err := h.exampleService.Create(ctx, &req)
	if err != nil {
		// Don't expose internal error details to client
		return c.Status(fiber.StatusInternalServerError).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", nil, requestID))
}
