package models

// AuthResponse represents the response payload for login operations
type AuthResponse struct {
	AccessToken         string `json:"access_token"`
	ExpiresIn          int64  `json:"expires_in"`
	RefreshToken       string `json:"refresh_token"`
	RefreshExpiresIn   int64  `json:"refresh_expires_in"`
}

// TokenRefreshResponse represents the response payload for token refresh
type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}