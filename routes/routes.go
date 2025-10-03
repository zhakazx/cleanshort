package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zhakazx/cleanshort/config"
	"github.com/zhakazx/cleanshort/controllers"
	"github.com/zhakazx/cleanshort/middleware"
	"github.com/zhakazx/cleanshort/services"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	authService := services.NewAuthService(db, cfg)
	linkService := services.NewLinkService(db, cfg)

	authController := controllers.NewAuthController(authService)
	linkController := controllers.NewLinkController(linkService)

	app.Get("/:shortCode",
		middleware.RedirectRateLimitMiddleware(cfg.RateLimitRedirect),
		linkController.RedirectLink,
	)

	// API v1 routes
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Use(middleware.AuthRateLimitMiddleware(cfg.RateLimitAuth))

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)
	auth.Post("/logout", authController.Logout)

	links := api.Group("/links")
	links.Use(middleware.AuthMiddleware(cfg))

	links.Post("/", linkController.CreateLink)
	links.Get("/", linkController.ListLinks)
	links.Get("/:id", linkController.GetLink)
	links.Patch("/:id", linkController.UpdateLink)
	links.Delete("/:id", linkController.DeleteLink)

	// Start cleanup goroutine for expired tokens
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Run daily
		defer ticker.Stop()

		for range ticker.C {
			if err := authService.CleanupExpiredTokens(); err != nil {
				continue
			}
		}
	}()
}
