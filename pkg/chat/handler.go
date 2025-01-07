package chat

import (
	"net/http"

	"StyleSwap/config"
	"StyleSwap/utils"

	"github.com/gorilla/websocket"
)

var (
	clientManager = NewClientManager()
	upgrader      = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Structure principale
type MessageConfig struct {
	*config.Config
}

// Nouvelle instance de MessageConfig
func New(configuration *config.Config) *MessageConfig {
	return &MessageConfig{configuration}
}

// HandleWebSocket gère la connexion WebSocket
func (config *MessageConfig) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Mettre à niveau la connexion HTTP -> WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err, "Erreur lors de la mise à niveau de la connexion HTTP -> WebSocket")
	defer conn.Close()

	// Authentifier l'utilisateur
	userID, err := AuthenticateUser(conn, r)
	utils.HandleError(err, "Erreur lors de l'authentification de l'utilisateur")

	// Ajouter le client
	clientManager.AddClient(userID, conn)
	defer clientManager.RemoveClient(userID)

	// Charger l'historique des messages
	config.GetConversation(r, userID)

	// Écouter les messages
	config.HandleMessage(userID, conn, r)
}
