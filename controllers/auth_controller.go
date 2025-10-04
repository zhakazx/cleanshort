package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zhakazx/cleanshort/models"
	"github.com/zhakazx/cleanshort/services"
	"github.com/zhakazx/cleanshort/utils"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req models.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid request body",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.HandleValidationError(c, err)
	}

	user, err := ac.authService.Register(&req)
	if err != nil {
		if strings.Contains(err.Error(), "already in use") {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "CONFLICT",
					Message:   "Email already in use",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to register user",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid request body",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.HandleValidationError(c, err)
	}

	authResponse, err := ac.authService.Login(&req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid credentials",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to authenticate user",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(authResponse)
}

func (ac *AuthController) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid request body",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.HandleValidationError(c, err)
	}

	tokenResponse, err := ac.authService.RefreshToken(&req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") || strings.Contains(err.Error(), "revoked") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid or expired refresh token",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to refresh token",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(tokenResponse)
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid request body",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.HandleValidationError(c, err)
	}

	if err := ac.authService.Logout(req.RefreshToken); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "NOT_FOUND",
					Message:   "Refresh token not found",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to logout user",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}