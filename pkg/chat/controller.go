package chat

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Structure pour stocker les connexions WebSocket
type Client struct {
	ID            string
	Conn          *websocket.Conn
	activeConvs   map[string]bool
	CurrentClient string
}

type MessageConfig struct {
	*config.Config
}

func New(configuration *config.Config) *MessageConfig {
	return &MessageConfig{configuration}
}

var (
	clients     = make(map[string]*Client) // Dictionnaire pour les clients connectés
	clientsLock sync.Mutex                 // Verrou pour la synchronisation des accès
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (config *MessageConfig) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Mettre à niveau la connexion HTTP vers WebSocket
	userID, ok := auth.GetUserIDFromContext(r.Context())

	if !ok {
		log.Println("ID utilisateur non trouvé")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connexion WebSocket établie")
	if err != nil {
		log.Println("Erreur WebSocket:", err)
		return
	}
	defer conn.Close()

	var reqAuth model.MessageRequest
	errRead := conn.ReadJSON(&reqAuth)

	// Valider la requête d'authentification
	if reqAuth.UserID == "" {
		log.Println("Champ manquant dans la requête")
		return
	}

	if errRead != nil {
		log.Printf("Erreur lors de la lecture de l'ID de l'utilisateur: %v\n", err)
		return
	}

	clientsLock.Lock()
	if _, exists := clients[userID]; exists {
		log.Println("Cet utilisateur est déjà connecté.")
		clientsLock.Unlock()
		return
	}

	clients[userID] = &Client{
		ID:          userID,
		Conn:        conn,
		activeConvs: make(map[string]bool), // Initialiser les conversations actives
	}
	clientsLock.Unlock()

	fmt.Printf("Client %s connecté\n", userID)

	// Récupérer les messages non livrés pour l'utilisateur
	config.GetConversation(userID, reqAuth.UserID)

	// Traitement des messages reçus
	for {
		var req model.MessageRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Erreur lors de la lecture du message:", err)
			break
		}

		// Changer de conversation
		if req.UserID != "" && req.UserID != reqAuth.UserID && req.UserID != userID {

			clientsLock.Lock()
			client := clients[userID]
			client.CurrentClient = req.UserID
			reqAuth.UserID = req.UserID // Met à jour le destinataire de la conversation
			clientsLock.Unlock()

			// Récupérer les anciens messages pour la nouvelle conversation
			config.GetConversation(userID, req.UserID)
		}

		// Sauvegarder et envoyer le message
		if req.Content == "" {
			continue
		}

		message := dbmodel.Messages{
			SenderID:   userID,
			ReceiverID: reqAuth.UserID,
			Content:    req.Content,
		}
		config.MessageRepository.Create(&message)

		// Vérifier si le destinataire est connecté
		clientsLock.Lock()
		destClient, ok := clients[reqAuth.UserID]
		clientsLock.Unlock()

		if ok {
			// Si le destinataire est connecté, envoyer le message
			reqAuth.Content = userID + ": " + reqAuth.Content
			err := destClient.Conn.WriteMessage(websocket.TextMessage, []byte(req.Content))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message au destinataire:", err)
			}
			message.Delivered = 0
		} else {
			log.Printf("Destinataire %s non trouvé, message sauvegardé\n", reqAuth.UserID)
			// Le message est marqué comme non livré, à envoyer lors de la reconnexion
			message.Delivered = 1
			config.MessageRepository.Save(&message)
		}
	}

	// Déconnecter le client lorsqu'il quitte
	clientsLock.Lock()
	delete(clients, userID)
	clientsLock.Unlock()
}