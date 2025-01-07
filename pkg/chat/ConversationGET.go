package chat

// import (
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// func (config *MessageConfig) GetConversation2(user string, client string) {

// 	if config.MessageRepository == nil {
// 		log.Fatal("Base de données non initialisée")
// 		return
// 	}

// 	messages := config.MessageRepository.GetConversation(user, client)

// 	clientsLock.Lock()
// 	client2, ok := clients[user]
// 	clientsLock.Unlock()

// 	if !ok {
// 		log.Printf("Client %s non trouvé\n", user)
// 		return
// 	}

// 	for _, message := range messages {
// 		if message.SenderID == user {
// 			message.Content = user + ": " + message.Content
// 		} else {
// 			message.Content = client + ": " + message.Content
// 		}
// 		err := client2.Conn.WriteMessage(websocket.TextMessage, []byte(message.Content))
// 		if err != nil {
// 			log.Println("Erreur lors de l'envoi du message au client:", err)
// 			continue
// 		}
// 	}
// }
