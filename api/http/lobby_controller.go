package http

/**
router.HandleFunc("/lobbies/{accessCode}", http.GetLobby).Methods("GET")
	router.HandleFunc("/lobbies/{accessCode}/join", http.JoinLobby).Methods("POST")
	router.HandleFunc("/lobbies", http.CreateLobby).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/select/player", http.SelectPlayer).Methods("POST")
	router.HandleFunc("/lobbies/{accessCode}/select/team", http.SelectTeam).Methods("POST")
*/

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paq-devs/paq-be-rpg/config"
	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

type ProfileRequest struct {
	ProfileId  string              `json:"profile_id"`
	Avatar     string              `json:"avatar"`
	Name       string              `json:"name"`
	HardSkills []profile.HardSkill `json:"hard_skills"`
	SoftSkills []profile.SoftSkill `json:"soft_skills"`
}

type SelectPlayerRequest struct {
	PlayerId string `json:"player_id"`
	LeaderId string `json:"leader_id"`
}

type SelectTeamRequest struct {
	TeamId   int    `json:"team_id"`
	LeaderId string `json:"leader_id"`
}

type LobbyCreateRequest struct {
	MasterName    string `json:"master_name"`
	MasterAvatar  string `json:"master_avatar"`
	LobbyName     string `json:"name"`
	MaxHardSkills int    `json:"max_hard_skills"`
	MaxSoftSkills int    `json:"max_soft_skills"`
}

// CreateLobby godoc
// @Summary Create a lobby
// @Description Create a lobby
// @Tags lobbies
// @Accept json
// @Produce json
// @Param request body LobbyCreateRequest true "Lobby request"
// @Success 200 {object} LobbyResponse
// @Router /lobbies [post]
func CreateLobby(w http.ResponseWriter, r *http.Request) {
	request := LobbyCreateRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby, err := config.GetModule().LobbyService.CreateLobby(r.Context(),
		profile.NewMaster(request.MasterName, request.MasterAvatar),
		request.LobbyName,
		request.MaxHardSkills,
		request.MaxSoftSkills)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// GetLobby godoc
// @Summary Get a lobby
// @Description Get a lobby by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode} [get]
func GetLobby(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]

	lobby, err := config.GetModule().LobbyService.GetLobby(r.Context(), accessCode)

	if (err != nil && err.Error() == "lobby_not_found") || lobby == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// JoinLobby godoc
// @Summary Join a lobby
// @Description Join a lobby by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/join [post]
func JoinLobby(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	request := ProfileRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby, err := config.GetModule().LobbyService.JoinLobby(r.Context(), accessCode, profile.NewPlayer(request.Name, request.Avatar, request.HardSkills, request.SoftSkills))

	if (err != nil && err.Error() == "lobby_not_found") || lobby == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// SelectPlayer godoc
// @Summary Select a player
// @Description Select a player by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/select/player [post]
func SelectPlayer(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	request := SelectPlayerRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby, err := config.GetModule().LobbyService.SelectPlayer(r.Context(), accessCode, profile.Profile{
		ID: request.LeaderId,
	}, request.PlayerId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if lobby == nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// SelectTeam godoc
// @Summary Select a team
// @Description Select a team by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/select/team [post]
func SelectTeam(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	request := SelectTeamRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby, err := config.GetModule().LobbyService.SelectTeam(r.Context(), accessCode, profile.Profile{
		ID: request.LeaderId,
	}, request.TeamId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if lobby == nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// JoinAsMentor godoc
// @Summary Join as mentor
// @Description Join as mentor by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/mentor [post]
func JoinMentor(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	request := ProfileRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobby, err := config.GetModule().LobbyService.JoinLobby(r.Context(), accessCode, profile.NewMentor(request.Name, request.Avatar))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if lobby == nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// CloseLobby godoc
// @Summary Close a lobby
// @Description Close a lobby by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/close [post]
func CloseLobby(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	lobby, err := config.GetModule().LobbyService.StartTeamCreation(r.Context(), accessCode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if lobby == nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}

// PromotePlayer godoc
// @Summary Promote a player
// @Description Promote a player by access code
// @Tags lobbies
// @Accept json
// @Produce json
// @Param accessCode path string true "Access code"
// @Param playerId path string true "Player id"
// @Success 200 {object} LobbyResponse
// @Failure 404 {object} ErrorResponse
// @Router /lobbies/{accessCode}/promote/{playerId} [post]
func PromotePlayer(w http.ResponseWriter, r *http.Request) {
	accessCode := mux.Vars(r)["accessCode"]
	playerId := mux.Vars(r)["playerId"]

	lobby, err := config.GetModule().LobbyService.PromoteLeader(r.Context(), accessCode, profile.Profile{
		ID: playerId,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if lobby == nil {
		http.Error(w, "lobby not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(lobby)
}
