package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index:idx_links_user"`
	ShortCode     string     `json:"short_code" gorm:"type:varchar(32);uniqueIndex;not null" validate:"required,min=4,max=32,alphanum"`
	TargetURL     string     `json:"target_url" gorm:"type:text;not null" validate:"required,url,max=2048"`
	Title         *string    `json:"title" gorm:"type:text"`
	IsActive      bool       `json:"is_active" gorm:"not null;default:true"`
	ClickCount    int64      `json:"click_count" gorm:"not null;default:0"`
	LastClickedAt *time.Time `json:"last_clicked_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;default:now()"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hook to generate UUID if not set
func (l *Link) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// LinkCreateRequest represents the request payload for creating a link
type LinkCreateRequest struct {
	TargetURL string  `json:"target_url" validate:"required,url,max=2048"`
	ShortCode *string `json:"short_code,omitempty" validate:"omitempty,min=4,max=32,alphanum"`
	Title     *string `json:"title,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// LinkUpdateRequest represents the request payload for updating a link
type LinkUpdateRequest struct {
	TargetURL *string `json:"target_url,omitempty" validate:"omitempty,url,max=2048"`
	Title     *string `json:"title,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// LinkResponse represents the response payload for link operations
type LinkResponse struct {
	ID            uuid.UUID  `json:"id"`
	ShortCode     string     `json:"short_code"`
	ShortURL      string     `json:"short_url"`
	TargetURL     string     `json:"target_url"`
	Title         *string    `json:"title"`
	IsActive      bool       `json:"is_active"`
	ClickCount    int64      `json:"click_count"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// LinkListResponse represents the response payload for listing links
type LinkListResponse struct {
	Links []LinkResponse `json:"links"`
	Total int64          `json:"total"`
	Limit int            `json:"limit"`
	Offset int           `json:"offset"`
}