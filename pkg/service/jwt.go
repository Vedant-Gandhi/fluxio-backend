package service

import (
	"fluxio-backend/pkg/common/schema"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
	logger schema.Logger
}

func NewJWTService(secret string, logger schema.Logger) *JWTService {
	return &JWTService{secret: secret, logger: logger}
}

func (s *JWTService) GenerateToken(payload model.JWTTokenClaims, exp time.Time) (token string, err error) {
	// Set default values for the payload
	if strings.EqualFold(payload.Sub, "") {
		payload.Sub = "user"
	}

	if strings.EqualFold(payload.UserID, "") {
		err = fluxerrors.ErrInvalidClaims
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": payload.UserID,
		"sub":     "user",
		"exp":     exp.Unix(),
		"iat":     time.Now().Unix(),
	})
	token, err = t.SignedString([]byte(s.secret))

	if err != nil {
		err = fluxerrors.ErrFailedToCreateToken
	}
	return
}

func (s *JWTService) ValidateToken(token string) (payload model.JWTTokenClaims, err error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fluxerrors.ErrInvalidToken
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		err = fluxerrors.ErrInvalidToken
	}

	claims := t.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(string)
	if strings.EqualFold(userID, "") {
		err = fluxerrors.ErrInvalidToken
	}

	payload.UserID = userID

	return
}
