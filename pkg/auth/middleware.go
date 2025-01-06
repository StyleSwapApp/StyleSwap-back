package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// AuthMiddleware protège les routes et extrait user_id depuis le token JWT
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer le header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Extraire le token JWT du header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// Parser le token JWT pour extraire user_id
			userID, err := ParseToken(secret, tokenString)
			if err != nil {
				log.Println("Token parsing error:", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Ajouter user_id au contexte
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext récupère user_id depuis le contexte
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}
