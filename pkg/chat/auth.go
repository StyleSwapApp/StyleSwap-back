package chat

import (
	"log"
	"net/http"
	"StyleSwap/pkg/auth"
	"github.com/gorilla/websocket"
)

// Authentifie un utilisateur via WebSocket
func AuthenticateUser(conn *websocket.Conn, r *http.Request) (string, error) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("ID utilisateur non trouvé")
		return "", nil
	}
	return userID, nil
}
