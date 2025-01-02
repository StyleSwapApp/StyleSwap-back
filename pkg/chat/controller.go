package chat

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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
	db *gorm.DB
)

type UserMessage struct {
    UserID string `json:"userID"`
}

func (config *MessageConfig) sendPendingMessages(userID string, conn *websocket.Conn) {
	// Récupérer les messages non livrés pour ce client
	
	Messages, err := config.MessageRepository.FindDelivered(userID)
	fmt.Println(Messages)
	if err != nil {
		log.Println("Erreur lors de la récupération des messages non livrés:", err)
		return
	}

	// Envoyer les messages non livrés au client
	for _, msg := range Messages {
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Notification: %s", msg.Content)))
		if err != nil {
			log.Println("Erreur lors de l'envoi de la notification de message:", err)
			break
		}
		// Marquer le message comme livré après l'envoi
		msg.Delivered = true
		db.Save(&msg)
	}
}

func (config *MessageConfig) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Mettre à niveau la connexion HTTP vers WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("Connexion WebSocket établie")
	if err != nil {
		log.Println("Erreur WebSocket:", err)
		return
	}
	defer conn.Close()

	var req UserMessage
	errRead := conn.ReadJSON(&req)

	if errRead != nil {
		log.Printf("Erreur lors de la lecture de l'ID de l'utilisateur: %v\n", err)
		return
	}

	clientsLock.Lock()
	if _, exists := clients[req.UserID]; exists {
		log.Println("Cet utilisateur est déjà connecté.")
		clientsLock.Unlock()
		return 
	}
	clients[req.UserID] = &Client{ID: req.UserID, Conn: conn}
	clientsLock.Unlock()


	// Ajouter le client à la map des connexions actives
	clientsLock.Lock()
	clients[req.UserID] = &Client{ID: req.UserID, Conn: conn}
	clientsLock.Unlock()

	fmt.Printf("Client %s connecté\n", req.UserID)

	config.sendPendingMessages(req.UserID, conn)

	if config.MessageRepository == nil {
		log.Fatal("Base de données non initialisée")
		return
	}	

	// Écouter les messages entrants et les rediriger vers le destinataire
	for {
		var msg string
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erreur WebSocket:", err)
			break
		}
		msg = string(msgBytes)

		// Le message doit être au format "receiverID:message"
		// Exemple : "client2:Hello"
		destID, content := parseMessage(msg)

		// Sauvegarder le message dans la base de données
		message := dbmodel.Messages{
			SenderID:   req.UserID,
			ReceiverID: destID,
			Content:    content,
		}
		config.MessageRepository.Create(&message)

		// Envoyer le message au destinataire si connecté
		clientsLock.Lock()
		destClient, ok := clients[destID]
		clientsLock.Unlock()

		if ok {
			// Si le destinataire est connecté, envoyer le message
			err := destClient.Conn.WriteMessage(websocket.TextMessage, []byte(content))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message au destinataire:", err)
			}
			message.Delivered = true
		} else {
			log.Printf("Destinataire %s non trouvé, message sauvegardé\n", destID)
			// Le message est marqué comme non livré, à envoyer lors de la reconnexion
			message.Delivered = false
			config.MessageRepository.Save(&message)
		}
	}

	// Déconnecter le client lorsqu'il quitte
	clientsLock.Lock()
	delete(clients, req.UserID)
	clientsLock.Unlock()
}

func parseMessage(msg string) (string, string) {
	// Extraire l'ID du destinataire et le contenu du message
	parts := strings.SplitN(msg, ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// func handleSendMessage(w http.ResponseWriter, r *http.Request) {
// 	// Lire les données de la requête
// 	var data struct {
// 		SenderID   string `json:"sender_id"`
// 		ReceiverID string `json:"receiver_id"`
// 		Content    string `json:"content"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		http.Error(w, "Erreur de lecture de la requête", http.StatusBadRequest)
// 		return
// 	}

// 	// Sauvegarder le message dans la base de données
// 	message := dbmodel.Message{
// 		SenderID:   data.SenderID,
// 		ReceiverID: data.ReceiverID,
// 		Content:    data.Content,
// 	}
// 	if err := db.Create(&message).Error; err != nil {
// 		http.Error(w, "Erreur lors de l'enregistrement du message", http.StatusInternalServerError)
// 		return
// 	}

// 	// Répondre au client
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Message envoyé avec succès")
// }
