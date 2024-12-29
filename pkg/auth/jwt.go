package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secret, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": name,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (string, error) {
	fmt.Println(tokenString, secret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	fmt.Println(interface{}(token))

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(ok)
	fmt.Println(claims)

	if ok && token.Valid {
		fmt.Println(claims["email"])
		fmt.Println(token.Valid)
		return claims["email"].(string), nil
	}
	return "", err
}
