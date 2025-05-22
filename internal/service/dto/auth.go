package dto

type RegistrationRequest struct {
	Login    string `json:"login" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Login      string `json:"login" validate:"required"`
	Password   string `json:"password" validate:"required"`
	DeviceCode string `json:"deviceCode" validate:"required"`
}
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
