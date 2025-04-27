package model

type JWTTokenClaims struct {
	UserID string `json:"user_id"`
	Sub    string `json:"sub"`
}
