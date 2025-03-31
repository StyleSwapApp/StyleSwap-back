package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims définit la structure des données dans le JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken génère un JWT avec user_id et une expiration
func GenerateToken(secret, userID string) (string, error) {
	// Création des claims avec user_id et exp
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}

	// Création du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signature du token avec le secret
	return token.SignedString([]byte(secret))
}

// ParseToken valide le token et extrait user_id
func ParseToken(secret, tokenString string) (string, error) {
	// Parser le token JWT
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Extraire les claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Retourner user_id
	if claims.UserID == "" {
		return "", errors.New("user_id not found in token")
	}

	return claims.UserID, nil
}
