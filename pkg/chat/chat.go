package chat

import "StyleSwap/config"

// Initialiser le gestionnaire de WebSocket
func InitializeChatHandler(cfg *config.Config) *MessageConfig {
	return New(cfg)
}
