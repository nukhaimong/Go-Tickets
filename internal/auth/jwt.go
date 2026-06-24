package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey         = "secret_key"
	defaultTokenDuration = 24 * time.Hour // one day
)

type JWTClaims struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, email string, name string) (string, error)
	ValidateToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretkey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string) JWTService {
	if secretKey == "" {
		secretKey = jwtSecretKey
	}
	return &jwtService{
		secretkey:     secretKey,
		tokenDuration: defaultTokenDuration,
	}
}

func (js *jwtService) GenerateToken(userId uint, email string, name string) (string, error) {
	// create claims
	claims := JWTClaims{
		UserId: userId,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gotickets",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // create token with claims

	tokenStr, err := token.SignedString([]byte(js.secretkey)) // sign token with secret key
	if err != nil {
		return "", nil
	}
	return tokenStr, nil
}

func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing mathod: %v", token.Header["alg"])
		}
		return []byte(js.secretkey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Token validation failed: %w", err)
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid token")
}
