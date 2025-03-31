package chat

import (
	"github.com/gorilla/websocket"
)

// Client représente un utilisateur connecté via WebSocket
type Client struct {
	ID            string
	Conn          *websocket.Conn
	CurrentClient string
}

// Envoyer un message au client
func (c *Client) SendMessage(message string) error {
	return c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// Fermer la connexion WebSocket proprement
func (c *Client) Close() error {
	return c.Conn.Close()
}
