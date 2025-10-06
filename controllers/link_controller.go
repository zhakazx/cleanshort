package controllers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zhakazx/cleanshort/models"
	"github.com/zhakazx/cleanshort/services"
	"github.com/zhakazx/cleanshort/utils"
)

type LinkController struct {
	linkService *services.LinkService
}

func NewLinkController(linkService *services.LinkService) *LinkController {
	return &LinkController{
		linkService: linkService,
	}
}

func (lc *LinkController) CreateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var req models.LinkCreateRequest
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

	link, err := lc.linkService.CreateLink(userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid short code") {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "VALIDATION_ERROR",
					Message:   "Invalid short code format",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		if strings.Contains(err.Error(), "reserved") {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "VALIDATION_ERROR",
					Message:   "Short code is reserved",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "CONFLICT",
					Message:   "Short code already exists",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to create link",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(link)
}

func (lc *LinkController) GetLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	linkIDStr := c.Params("id")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid link ID",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	link, err := lc.linkService.GetLink(userID, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "LINK_NOT_FOUND",
					Message:   "Short link not found",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to retrieve link",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(link)
}

func (lc *LinkController) UpdateLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	linkIDStr := c.Params("id")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid link ID",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	var req models.LinkUpdateRequest
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

	link, err := lc.linkService.UpdateLink(userID, linkID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "LINK_NOT_FOUND",
					Message:   "Short link not found",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to update link",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(link)
}

func (lc *LinkController) DeleteLink(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	linkIDStr := c.Params("id")
	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid link ID",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	err = lc.linkService.DeleteLink(userID, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "LINK_NOT_FOUND",
					Message:   "Short link not found",
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to delete link",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (lc *LinkController) ListLinks(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Parse query parameters
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")
	query := c.Query("query", "")
	activeStr := c.Query("active", "")
	sortBy := c.Query("sort_by", "created_at")
	orderBy := c.Query("order_by", "desc")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var active *bool
	if activeStr != "" {
		switch activeStr {
		case "true":
			activeVal := true
			active = &activeVal
		case "false":
			activeVal := false
			active = &activeVal
		}
	}

	// Validate sort_by parameter
	validSortFields := map[string]bool{
		"created_at":       true,
		"updated_at":       true,
		"title":           true,
		"short_code":      true,
		"click_count":     true,
		"last_clicked_at": true,
	}

	if !validSortFields[sortBy] {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid sort_by field. Allowed values: created_at, updated_at, title, short_code, click_count, last_clicked_at",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	// Validate order_by parameter
	if orderBy != "asc" && orderBy != "desc" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "VALIDATION_ERROR",
				Message:   "Invalid order_by value. Allowed values: asc, desc",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	links, err := lc.linkService.ListLinks(userID, limit, offset, query, active, sortBy, orderBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "INTERNAL_ERROR",
				Message:   "Failed to retrieve links",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(links)
}

func (lc *LinkController) RedirectLink(c *fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	link, err := lc.linkService.GetLinkByShortCode(shortCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "LINK_NOT_FOUND",
				Message:   "Short link not found",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	if !link.IsActive {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error: models.ErrorDetail{
				Code:      "LINK_NOT_FOUND",
				Message:   "Short link not found",
				RequestID: c.Locals("requestid").(string),
			},
		})
	}

	// Record click (async)
	go func() {
		lc.linkService.RecordClick(shortCode)
	}()

	return c.Redirect(link.TargetURL, fiber.StatusFound)
}
