package helpers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"time"
)

const (
	TokenExpiration = 30 * time.Minute
)

var signingKey = []byte(os.Getenv("SECRET_KEY"))

//func GenerateToken(id uuid.UUID, username string) (string, error) {
//	return "", nil
//}

//func GenerateToken(userID uuid.UUID, email string) (string, error) {
//	var signingKey = []byte(os.Getenv("SECRET_KEY"))
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)
//
//	claims["email"] = email
//	claims["user_id"] = userID
//	claims["exp"] = time.Now().Add(TokenExpiration).Unix()
//
//	tokenString, err := token.SignedString(signingKey)
//	if err != nil {
//		fmt.Printf("error generating token: %s", err.Error())
//		return "", err
//	}
//
//	return tokenString, nil
//}

type Claims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(id uuid.UUID, email string) (string, error) {

	claims := Claims{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // expire in 1h
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "scrapy",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign with secret
	return token.SignedString(signingKey)
}

func DecodeToken(tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
