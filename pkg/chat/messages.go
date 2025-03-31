package chat

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// HandleMessage gère les messages WebSocket

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


// DeliverMessage envoie un message à un client WebSocket ou le sauvegarde en base de données si le client n'est pas connecté

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
		utils.HandleError(err, "Erreur lors de l'envoi du message au client")
		delivered = 0
	} else {
		delivered = 1
	}

	config.AjouterBDD(senderID, receiverID, content, delivered)
}


// AjouterBDD ajoute un message à la base de données

func (config *MessageConfig) AjouterBDD(SenderID string, ReceiverID string, content string, delivered int) {
	message := dbmodel.Messages{
		SenderID:   SenderID,
		ReceiverID: ReceiverID,
		Content:    content,
		Delivered:  delivered,
	}
	config.MessageRepository.Create(&message)
}
