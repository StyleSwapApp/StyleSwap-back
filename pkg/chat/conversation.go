package chat

import (
	"StyleSwap/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// GetConversation récupère et envoie les messages historiques
func (config *MessageConfig) GetConversation(r *http.Request, user string) {
	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}

	// Récupérer l'ID du client
	IDclient := chi.URLParam(r, "idclient")
	clientID, err := strconv.Atoi(IDclient)
	utils.HandleError(err, "Erreur lors de la conversion de l'ID du client en entier")

	client, err := config.UserRepository.FindByID(clientID)
	utils.HandleError(err, "Erreur lors de la recherche du client par ID")

	//Être sûr que le client n'est pas le même que l'utilisateur
	if client.Pseudo == user {
		log.Printf("Impossible de démarrer une conversation avec soi-même\n")
		return
	}

	// Récupérer les messages depuis la base de données
	messages := config.MessageRepository.GetConversation(user, client.Pseudo)

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
			formattedContent = client.Pseudo + ": " + message.Content
		}

		err := clientInstance.Conn.WriteMessage(websocket.TextMessage, []byte(formattedContent))
		utils.HandleError(err, "Erreur lors de l'envoi du message au client")
	}
}
