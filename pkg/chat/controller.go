package chat

import (
	"StyleSwap/config"
	"StyleSwap/pkg/model"
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
	userID, conn := config.init(w, r)
	defer conn.Close()

	var reqAuth model.MessageRequest
	errRead := conn.ReadJSON(&reqAuth)

	// Valider la requête d'authentification
	if reqAuth.UserID == "" {
		log.Println("Champ manquant dans la requête, dites moi à qui vous voulez envoyer le message")
		return
	}

	if errRead != nil {
		log.Printf("Erreur lors de la lecture de l'ID de l'utilisateur: %v\n", errRead)
		return
	}

	// Vérifier si l'utilisateur est connecté
	nouvelleConnexion(userID, conn)

	//Envoyer le message si le champ Content n'est pas vide
	if reqAuth.Content != "" {
		delivered := 1
		config.AjouterBDD(userID, reqAuth.UserID, reqAuth.Content, delivered)
	}
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

		// Vérifier si le message est vide
		errContent := req.Bind(r)
		if errContent != nil {
			continue
		}

		// Vérifier si le destinataire est connecté
		clientsLock.Lock()
		destClient, ok := clients[reqAuth.UserID]
		clientsLock.Unlock()

		// Envoyer le message au destinataire
		var delivered int
		if ok { // Si le destinataire est connecté
			messageEnvoye := userID + ": " + req.Content
			err := destClient.Conn.WriteMessage(websocket.TextMessage, []byte(messageEnvoye))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message au destinataire:", err)
			}
			delivered = 0
		} else { // Si le destinataire n'est pas connecté
			log.Printf("Destinataire %s non trouvé, message sauvegardé\n", reqAuth.UserID)
			delivered = 1
		}
		//ajouter le message à la base de données
		config.AjouterBDD(userID, reqAuth.UserID, req.Content, delivered)
	}

	// Déconnecter le client lorsqu'il quitte
	clientsLock.Lock()
	delete(clients, userID)
	clientsLock.Unlock()
}
