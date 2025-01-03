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
	ID   string
	Conn *websocket.Conn
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
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connexion WebSocket établie")
	if err != nil {
		log.Println("Erreur WebSocket:", err)
		return
	}
	defer conn.Close()

	var reqAuth model.AuthedRequest
	errRead := conn.ReadJSON(&reqAuth)

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
	clients[reqAuth.UserID] = &Client{ID: reqAuth.UserID, Conn: conn}
	clientsLock.Unlock()

	fmt.Printf("Client %s connecté\n", reqAuth.UserID)

	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}

	for {
		var req model.MessageRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			log.Println("Erreur lors de la lecture du message:", err)
			break
		}

		if req.Bind(r) != nil {
			log.Println("Champ manquant dans la requête")
			continue
		}

		message := dbmodel.Messages{
			SenderID:   reqAuth.UserID,
			ReceiverID: req.ReceiverID,
			Content:    req.Content,
		}
		config.MessageRepository.Create(&message)
		
		// Vérifier si le destinataire est connecté
		clientsLock.Lock()
		destClient, ok := clients[req.ReceiverID]
		clientsLock.Unlock()

		if ok {
			// Si le destinataire est connecté, envoyer le message
			err := destClient.Conn.WriteMessage(websocket.TextMessage, []byte(req.Content))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message au destinataire:", err)
			}
			message.Delivered = 0
		} else {
			log.Printf("Destinataire %s non trouvé, message sauvegardé\n", req.ReceiverID)
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

// func parseMessage(msg string) (string, string) {
// 	// Extraire l'ID du destinataire et le contenu du message
// 	if !isValidJSON(msg) {
// 		parts := strings.SplitN(msg, ":", 2)
// 		if len(parts) != 2 {
// 			return "", ""
// 		}
// 		return parts[0], parts[1]
// 	}
// 	return "", ""
// }

func (config *MessageConfig) SendMessagesToClient(userID string) {

	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}

	// Récupérer les messages non livrés pour l'utilisateur
	messages := config.MessageRepository.GetUndeliveredMessages(userID)
	if len(messages) == 0 {
		log.Printf("Aucun message non livré pour l'utilisateur %s\n", userID)
		return
	}

	clientsLock.Lock()
	client, ok := clients[userID]
	clientsLock.Unlock()

	if !ok {
		log.Printf("Client %s non trouvé\n", userID)
		return
	}

	for _, message := range messages {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message.Content))
		if err != nil {
			log.Println("Erreur lors de l'envoi du message au client:", err)
			continue
		}
		message.Delivered = 0
		config.MessageRepository.Save(&message)
	}
}
