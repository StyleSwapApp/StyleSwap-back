package chat

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (config *MessageConfig) init(w http.ResponseWriter, r *http.Request) (string, *websocket.Conn) {
	// Récupérer l'ID de l'utilisateur à partir du contexte
	userID, ok := auth.GetUserIDFromContext(r.Context())

	if !ok {
		log.Println("ID utilisateur non trouvé")
		return "", nil
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connexion WebSocket établie")
	if err != nil {
		log.Println("Erreur WebSocket:", err)
		return "", nil
	}

	return userID, conn
}

func (config *MessageConfig) AjouterBDD(SenderID string, ReceiverID string, content string, delivered int) {
	message := dbmodel.Messages{
		SenderID:   SenderID,
		ReceiverID: ReceiverID,
		Content:    content,
		Delivered:  delivered,
	}
	config.MessageRepository.Create(&message)
}