package internal

import (
	"errors"
	"time"

	"github.com/JonasLindermayr/FileBeam/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var secretKey = []byte("b4b2daceb20249ebca2a2b7750b6eba7d3ad2fd45fde1dff7ff9e21172637467")

func GenerateToken(userId uuid.UUID, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId.String(),
		"exp":    expiresAt.Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func DecodeToken(tokenString string) (types.JWT, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return types.JWT{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdStr, ok := claims["userId"].(string)
		if !ok {
			return types.JWT{}, errors.New("userId claim missing or invalid")
		}
		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			return types.JWT{}, errors.New("invalid userId format")
		}
		return types.JWT{UUID: userId}, nil
	} else {
		return types.JWT{}, errors.New("invalid token claims")
	}
}
