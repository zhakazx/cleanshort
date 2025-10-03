package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zhakazx/cleanshort/config"
	"github.com/zhakazx/cleanshort/models"
	"github.com/zhakazx/cleanshort/utils"
	"gorm.io/gorm"
)

type LinkService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewLinkService(db *gorm.DB, cfg *config.Config) *LinkService {
	return &LinkService{
		db:  db,
		cfg: cfg,
	}
}

func (s *LinkService) CreateLink(userID uuid.UUID, req *models.LinkCreateRequest) (*models.LinkResponse, error) {
	var shortCode string
	var err error

	if req.ShortCode != nil && *req.ShortCode != "" {
		shortCode = strings.TrimSpace(*req.ShortCode)

		if !utils.IsValidShortCode(shortCode) {
			return nil, errors.New("invalid short code format")
		}

		if utils.IsReservedShortCode(shortCode) {
			return nil, errors.New("short code is reserved")
		}

		var existingLink models.Link
		if err := s.db.Where("short_code = ?", shortCode).First(&existingLink).Error; err == nil {
			return nil, errors.New("short code already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		shortCode, err = s.generateUniqueShortCode()
		if err != nil {
			return nil, err
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	link := models.Link{
		UserID:    userID,
		ShortCode: shortCode,
		TargetURL: req.TargetURL,
		Title:     req.Title,
		IsActive:  isActive,
	}

	if err := s.db.Create(&link).Error; err != nil {
		return nil, err
	}

	return s.linkToResponse(&link), nil
}

func (s *LinkService) GetLink(userID, linkID uuid.UUID) (*models.LinkResponse, error) {
	var link models.Link
	if err := s.db.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("link not found")
		}
		return nil, err
	}

	return s.linkToResponse(&link), nil
}

func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	if err := s.db.Where("short_code = ?", shortCode).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("link not found")
		}
		return nil, err
	}

	return &link, nil
}

func (s *LinkService) UpdateLink(userID, linkID uuid.UUID, req *models.LinkUpdateRequest) (*models.LinkResponse, error) {
	var link models.Link
	if err := s.db.Where("id = ? AND user_id = ?", linkID, userID).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("link not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.TargetURL != nil {
		updates["target_url"] = *req.TargetURL
	}

	if req.Title != nil {
		updates["title"] = *req.Title
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.db.Model(&link).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := s.db.Where("id = ?", linkID).First(&link).Error; err != nil {
		return nil, err
	}

	return s.linkToResponse(&link), nil
}

func (s *LinkService) DeleteLink(userID, linkID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", linkID, userID).Delete(&models.Link{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("link not found")
	}

	return nil
}

func (s *LinkService) ListLinks(userID uuid.UUID, limit, offset int, query string, active *bool) (*models.LinkListResponse, error) {
	var links []models.Link
	var total int64

	db := s.db.Model(&models.Link{}).Where("user_id = ?", userID)

	if active != nil {
		db = db.Where("is_active = ?", *active)
	}

	if query != "" {
		searchPattern := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(short_code) LIKE ? OR LOWER(title) LIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&links).Error; err != nil {
		return nil, err
	}

	linkResponses := make([]models.LinkResponse, len(links))
	for i, link := range links {
		linkResponses[i] = *s.linkToResponse(&link)
	}

	return &models.LinkListResponse{
		Links:  linkResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *LinkService) RecordClick(shortCode string) error {
	now := time.Now()
	return s.db.Model(&models.Link{}).
		Where("short_code = ?", shortCode).
		Updates(map[string]interface{}{
			"click_count":     gorm.Expr("click_count + 1"),
			"last_clicked_at": now,
		}).Error
}

func (s *LinkService) generateUniqueShortCode() (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		shortCode, err := utils.GenerateShortCode(8)
		if err != nil {
			return "", err
		}

		if utils.IsReservedShortCode(shortCode) {
			continue
		}

		var existingLink models.Link
		if err := s.db.Where("short_code = ?", shortCode).First(&existingLink).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return shortCode, nil
			}
			return "", err
		}
	}

	return "", errors.New("failed to generate unique short code")
}

func (s *LinkService) linkToResponse(link *models.Link) *models.LinkResponse {
	return &models.LinkResponse{
		ID:            link.ID,
		ShortCode:     link.ShortCode,
		ShortURL:      fmt.Sprintf("%s/%s", s.cfg.BaseURL, link.ShortCode),
		TargetURL:     link.TargetURL,
		Title:         link.Title,
		IsActive:      link.IsActive,
		ClickCount:    link.ClickCount,
		LastClickedAt: link.LastClickedAt,
		CreatedAt:     link.CreatedAt,
		UpdatedAt:     link.UpdatedAt,
	}
}
