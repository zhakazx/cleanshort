package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zhakazx/cleanshort/config"
	"github.com/zhakazx/cleanshort/models"
	"github.com/zhakazx/cleanshort/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

func (s *AuthService) Register(req *models.UserRegisterRequest) (*models.UserResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Normalize email
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Store refresh token
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshTTL),
	}

	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:        accessToken,
		ExpiresIn:         int64(s.cfg.JWTAccessTTL.Seconds()),
		RefreshToken:       refreshToken,
		RefreshExpiresIn:  int64(s.cfg.JWTRefreshTTL.Seconds()),
	}, nil
}

func (s *AuthService) RefreshToken(req *models.RefreshTokenRequest) (*models.TokenRefreshResponse, error) {
	hasher := sha256.New()
	hasher.Write([]byte(req.RefreshToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	var refreshToken models.RefreshToken
	if err := s.db.Where("token_hash = ?", tokenHash).Preload("User").First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid refresh token")
		}
		return nil, err
	}

	if !refreshToken.IsValid() {
		return nil, errors.New("refresh token is expired or revoked")
	}

	accessToken, err := utils.GenerateAccessToken(refreshToken.User.ID, refreshToken.User.Email, s.cfg)
	if err != nil {
		return nil, err
	}

	return &models.TokenRefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(s.cfg.JWTAccessTTL.Seconds()),
	}, nil
}

func (s *AuthService) Logout(refreshTokenString string) error {
	hasher := sha256.New()
	hasher.Write([]byte(refreshTokenString))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	result := s.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

func (s *AuthService) RevokeAllUserTokens(userID uuid.UUID) error {
	return s.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (s *AuthService) CleanupExpiredTokens() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}