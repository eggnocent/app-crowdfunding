package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY = []byte("CrOwD_FunD_S3crEt_kEY")

type Service interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (js *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *JWTService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(encodedToken, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
