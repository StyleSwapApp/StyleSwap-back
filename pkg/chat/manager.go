package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ClientManager gère les clients connectés
type ClientManager struct {
	Clients map[string]*Client
	Lock    sync.Mutex
}

// Crée une nouvelle instance de ClientManager
func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients: make(map[string]*Client),
	}
}

// Ajouter un client
func (cm *ClientManager) AddClient(id string, conn *websocket.Conn) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	cm.Clients[id] = &Client{ID: id, Conn: conn}
}

// Supprimer un client
func (cm *ClientManager) RemoveClient(id string) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	delete(cm.Clients, id)
}

// Récupérer un client
func (cm *ClientManager) GetClient(id string) (*Client, bool) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	client, ok := cm.Clients[id]
	return client, ok
}
