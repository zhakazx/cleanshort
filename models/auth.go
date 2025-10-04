package models

type AuthResponse struct {
	AccessToken         string `json:"access_token"`
	ExpiresIn          int64  `json:"expires_in"`
	RefreshToken       string `json:"refresh_token"`
	RefreshExpiresIn   int64  `json:"refresh_expires_in"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}