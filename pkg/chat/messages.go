package chat

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Gère les messages entrants
func (config *MessageConfig) HandleMessage(userID string, conn *websocket.Conn, r *http.Request) {
	for {
		var req model.MessageRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Erreur lors de la lecture du message:", err)
			break
		}

		// Envoyer le message
		config.DeliverMessage(userID, req.Content, r)
	}
}

func (config *MessageConfig) DeliverMessage(receiverID, content string, r *http.Request) {
	senderID, okUser := auth.GetUserIDFromContext(r.Context())
	if !okUser {

		return
	}
	destClient, ok := clientManager.GetClient(receiverID)

	var delivered int
	if ok {
		message := senderID + ": " + content
		err := destClient.SendMessage(message)
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		delivered = 0
	} else {
		log.Printf("Destinataire %s non trouvé, message sauvegardé\n", receiverID)
		delivered = 1
	}

	config.AjouterBDD(senderID, receiverID, content, delivered)
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
