package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"StyleSwap/config"
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
	if err != nil {
		log.Println("Échec de l'upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	// Authentifier l'utilisateur
	userID, err := AuthenticateUser(conn, r)
	if err != nil {
		log.Println(err)
		return
	}

	// Ajouter le client
	clientManager.AddClient(userID, conn)
	defer clientManager.RemoveClient(userID)

	// Charger l'historique des messages
	config.GetConversation(r, userID)

	// Écouter les messages
	config.HandleMessage(userID, conn, r)
}
