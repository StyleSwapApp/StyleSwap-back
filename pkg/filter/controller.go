package filter

import (
	"net/http"
)

type FilterCriteria struct {
	Pseudo  string
	Couleur string
	Marque  string
	Taille  string
}

// ParseFilterCriteria extrait les filtres depuis une requête HTTP.
func ParseFilterCriteria(r *http.Request) *FilterCriteria {
	return &FilterCriteria{
		Pseudo:  r.URL.Query().Get("pseudo"),
		Couleur: r.URL.Query().Get("color"),
		Marque:  r.URL.Query().Get("brand"),
		Taille:  r.URL.Query().Get("size"),
	}
}