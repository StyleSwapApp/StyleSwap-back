package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// GetConversation récupère et envoie les messages historiques
func (config *MessageConfig) GetConversation(user string, client string) {
	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}

	// Récupérer les messages depuis la base de données
	messages := config.MessageRepository.GetConversation(user, client)

	// Récupérer le client depuis le gestionnaire
	clientInstance, ok := clientManager.GetClient(user)
	if !ok {
		log.Printf("Client %s non trouvé\n", user)
		return
	}

	// Envoyer les messages au client
	for _, message := range messages {
		var formattedContent string
		if message.SenderID == user {
			formattedContent = user + ": " + message.Content
		} else {
			formattedContent = client + ": " + message.Content
		}

		err := clientInstance.Conn.WriteMessage(websocket.TextMessage, []byte(formattedContent))
		if err != nil {
			log.Println("Erreur lors de l'envoi du message au client:", err)
		}
	}
}
