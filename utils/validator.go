package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/zhakazx/cleanshort/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns formatted error response
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func HandleValidationError(c *fiber.Ctx, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, e.Field()+" is required")
			case "email":
				errorMessages = append(errorMessages, e.Field()+" must be a valid email")
			case "min":
				errorMessages = append(errorMessages, e.Field()+" must be at least "+e.Param()+" characters")
			case "max":
				errorMessages = append(errorMessages, e.Field()+" must be at most "+e.Param()+" characters")
			case "url":
				errorMessages = append(errorMessages, e.Field()+" must be a valid URL")
			case "alphanum":
				errorMessages = append(errorMessages, e.Field()+" must contain only alphanumeric characters")
			default:
				errorMessages = append(errorMessages, e.Field()+" is invalid")
			}
		}

		message := "Validation failed"
		if len(errorMessages) > 0 {
			message = errorMessages[0]
		}

		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   message,
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
		Error: models.ErrorDetail{
			Code:      "VALIDATION_ERROR",
			Message:   err.Error(),
			RequestID: c.Locals("requestid").(string),
		},
	})
}