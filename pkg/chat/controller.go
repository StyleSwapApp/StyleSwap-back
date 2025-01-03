package chat

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
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


// connexion WebSocket sous cette forme:
// {
//     "userID":"Simon",
//     "clientID":"swagger"
// }

// message WebSocket sous cette forme:
// {
//     "content":"Bonjour"
// }

//changement client WebSocket sous cette forme:
// {
//     "userID":"Simon",
//     "content":"bonjour Simon"
// }

func (config *MessageConfig) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Mettre à niveau la connexion HTTP vers WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connexion WebSocket établie")
	if err != nil {
		log.Println("Erreur WebSocket:", err)
		return
	}
	defer conn.Close()

	var reqAuth model.AuthedRequest
	errRead := conn.ReadJSON(&reqAuth)

	// Valider la requête d'authentification
	if reqAuth.Bind(r) != nil {
		log.Println("Champ manquant dans la requête")
		return
	}

	if errRead != nil {
		log.Printf("Erreur lors de la lecture de l'ID de l'utilisateur: %v\n", err)
		return
	}

	clientsLock.Lock()
	if _, exists := clients[reqAuth.UserID]; exists {
		log.Println("Cet utilisateur est déjà connecté.")
		clientsLock.Unlock()
		return
	}

	clients[reqAuth.UserID] = &Client{
		ID:          reqAuth.UserID,
		Conn:        conn,
		activeConvs: make(map[string]bool), // Initialiser les conversations actives
	}
	clientsLock.Unlock()

	fmt.Printf("Client %s connecté\n", reqAuth.UserID)

	// Récupérer les messages non livrés pour l'utilisateur
	config.GetConversation(reqAuth.UserID, reqAuth.ClientID)

	// Traitement des messages reçus
	for {
		var req model.MessageRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Erreur lors de la lecture du message:", err)
			break
		}

		// Changer de conversation
		if req.UserID != "" && req.UserID != reqAuth.UserID {

			clientsLock.Lock()
			client := clients[reqAuth.UserID]
			client.CurrentClient = req.UserID
			reqAuth.ClientID = req.UserID // Met à jour le destinataire de la conversation
			clientsLock.Unlock()

			// Récupérer les anciens messages pour la nouvelle conversation
			config.GetConversation(reqAuth.UserID, req.UserID)
		}

		// Sauvegarder et envoyer le message
		if req.Content == "" {
			continue
		}

		message := dbmodel.Messages{
			SenderID:   reqAuth.UserID,
			ReceiverID: reqAuth.ClientID,
			Content:    req.Content,
		}
		config.MessageRepository.Create(&message)

		// Vérifier si le destinataire est connecté
		clientsLock.Lock()
		destClient, ok := clients[reqAuth.ClientID]
		clientsLock.Unlock()

		if ok {
			// Si le destinataire est connecté, envoyer le message
			req.Content = reqAuth.ClientID + ": " + req.Content
			err := destClient.Conn.WriteMessage(websocket.TextMessage, []byte(req.Content))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message au destinataire:", err)
			}
			message.Delivered = 0
		} else {
			log.Printf("Destinataire %s non trouvé, message sauvegardé\n", reqAuth.ClientID)
			// Le message est marqué comme non livré, à envoyer lors de la reconnexion
			message.Delivered = 1
			config.MessageRepository.Save(&message)
		}
	}

	// Déconnecter le client lorsqu'il quitte
	clientsLock.Lock()
	delete(clients, reqAuth.UserID)
	clientsLock.Unlock()
}

func (config *MessageConfig) GetConversation(user string, client string) {

	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}

	messages := config.MessageRepository.GetConversation(user, client)

	clientsLock.Lock()
	client2, ok := clients[user]
	clientsLock.Unlock()

	if !ok {
		log.Printf("Client %s non trouvé\n", user)
		return
	}

	for _, message := range messages {
		if message.SenderID == user {
			message.Content = user + ": " + message.Content
		} else {
			message.Content = client + ": " + message.Content
		}
		err := client2.Conn.WriteMessage(websocket.TextMessage, []byte(message.Content))
		if err != nil {
			log.Println("Erreur lors de l'envoi du message au client:", err)
			continue
		}
	}
}
