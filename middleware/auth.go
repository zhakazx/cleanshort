package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zhakazx/cleanshort/config"
	"github.com/zhakazx/cleanshort/models"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Authorization header is required",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid authorization header format",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid or expired token",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid token claims",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		// Parse user ID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "UNAUTHORIZED",
					Message:   "Invalid user ID in token",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		// Set user context
		c.Locals("userID", userID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}