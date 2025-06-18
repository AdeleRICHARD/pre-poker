package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// Structures de base
type Ticket struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Vote struct {
	TicketID string `json:"ticketId"`
	User     string `json:"user"`
	Points   int    `json:"points"`
}

// Stockage en mémoire
var (
	votes      = make(map[string][]Vote) // ticketID -> votes
	votesMutex sync.Mutex
)

// Handler pour récupérer les tickets Jira (POC: tickets mockés)
func getTicketsHandler(w http.ResponseWriter, r *http.Request) {
	// Pour le POC, on retourne des tickets mockés
	tickets := []Ticket{
		{ID: "JIRA-1", Title: "Corriger le bug d'affichage"},
		{ID: "JIRA-2", Title: "Ajouter la page de login"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

// Handler pour enregistrer un vote
func voteHandler(w http.ResponseWriter, r *http.Request) {
	var v Vote
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lecture body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		http.Error(w, "Erreur JSON", http.StatusBadRequest)
		return
	}

	votesMutex.Lock()
	defer votesMutex.Unlock()
	votes[v.TicketID] = append(votes[v.TicketID], v)

	w.WriteHeader(http.StatusCreated)
}

// Handler pour récupérer les votes d'un ticket (pour debug/POC)
func getVotesHandler(w http.ResponseWriter, r *http.Request) {
	ticketID := r.URL.Query().Get("ticketId")
	votesMutex.Lock()
	defer votesMutex.Unlock()
	json.NewEncoder(w).Encode(votes[ticketID])
}

func main() {
	http.HandleFunc("/tickets", getTicketsHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/votes", getVotesHandler) // ?ticketId=JIRA-1

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serveur démarré sur http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
