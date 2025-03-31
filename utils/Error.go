package utils

import (
	"log"
)

// HandleError gère les erreurs et retourne un booléen pour savoir si une erreur a été rencontrée
func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v\n", message, err)
		return 
	}
}