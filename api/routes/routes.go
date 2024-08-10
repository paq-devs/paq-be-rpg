package routes

import (
	net_http "net/http"

	"github.com/gorilla/mux"
	"github.com/paq-devs/paq-be-rpg/api/http"
)

func contentTypeMiddleware(next net_http.Handler) net_http.Handler {
	return net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(contentTypeMiddleware)

	router.HandleFunc("/lobbies", http.CreateLobby).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}", http.GetLobby).Methods("GET")
	router.HandleFunc("/lobbies/{accessCode}/join", http.JoinLobby).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/join/mentor", http.JoinMentor).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/select/player", http.SelectPlayer).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/select/team", http.SelectTeam).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/close", http.CloseLobby).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/promote/{playerId}", http.PromotePlayer).Methods("POST")

	return router
}
