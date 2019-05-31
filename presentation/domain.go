package presentation

import (
	"encoding/json"
	"net/http"

	"../business"
	"github.com/go-chi/chi"
)

// GetDomains - Get the domains with their servers ordered by id (hostname)
func GetDomains(w http.ResponseWriter, r *http.Request) {
	var domains = business.GetDomains()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&domains)
}

// CreateDomain - Creates and returns a domain given the id (hostname)
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	domain, err := business.CreateDomain(id)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(&err)
	} else {
		json.NewEncoder(w).Encode(&domain)
	}

}
