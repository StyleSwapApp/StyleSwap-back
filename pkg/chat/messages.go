package chat

import (
	"StyleSwap/pkg/model"
	"log"

	"github.com/gorilla/websocket"
)

// Gère les messages entrants
func (config *MessageConfig) HandleMessage(userID string, conn *websocket.Conn) {
	for {
		var req model.MessageRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Erreur lors de la lecture du message:", err)
			break
		}

		// Changement de conversation
		if req.UserID != "" {
			config.SwitchConversation(userID, req.UserID)
		}

		// Envoyer le message
		config.DeliverMessage(userID, req.UserID, req.Content)
	}
}

// Changer de conversation
func (config *MessageConfig) SwitchConversation(userID, newUserID string) {
	client, ok := clientManager.GetClient(userID)
	if ok {
		client.CurrentClient = newUserID
		config.GetConversation(userID, newUserID)
	}
}

// Envoyer un message
func (config *MessageConfig) DeliverMessage(senderID, receiverID, content string) {
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
