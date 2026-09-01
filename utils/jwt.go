package utils

import (
	"movie-review-api/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	SecretKey      string
	ExpirationTime int64
}

func NewJWTUtil(secretKey string, expirationMs int64) *JWTUtil {
	return &JWTUtil{
		SecretKey:      secretKey,
		ExpirationTime: expirationMs,
	}
}

func (j *JWTUtil) GenerateToken(user *models.User) (string, error) {
	expiryTime := time.Now().Add(time.Duration(j.ExpirationTime) * time.Millisecond)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiryTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTUtil) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
