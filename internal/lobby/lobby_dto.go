package lobby

import "github.com/paq-devs/paq-be-rpg/internal/profile"

type ProfileResponse struct {
	ID            string              `json:"id"`
	Avatar        string              `json:"avatar"`
	Name          string              `json:"name"`
	Role          profile.Role        `json:"role"`
	HardSkills    []profile.HardSkill `json:"hard_skills"`
	SoftSkills    []profile.SoftSkill `json:"soft_skills"`
	JoinTimestamp int64               `json:"join_timestamp"`
}

type TeamResponse struct {
	ID      int               `json:"id"`
	Leader  ProfileResponse   `json:"leader"`
	Players []ProfileResponse `json:"players"`
	Mentor  ProfileResponse   `json:"mentor"`
}

type LobbyResponse struct {
	AccessCode    string            `json:"access_code"`
	Name          string            `json:"name"`
	MaxHardSkills int               `json:"max_hard_skills"`
	MaxSoftSkills int               `json:"max_soft_skills"`
	Status        LobbyStatus       `json:"status"`
	Players       []ProfileResponse `json:"players"`
	Mentors       []ProfileResponse `json:"mentors"`
	Master        ProfileResponse   `json:"master"`
	Teams         []TeamResponse    `json:"teams"`
	ChooseControl *ChooseControl    `json:"choose_control"`
}

func ResponseFromProfile(p *profile.Profile) ProfileResponse {
	return ProfileResponse{
		ID:            p.ID,
		Avatar:        p.Avatar,
		Name:          p.Name,
		Role:          p.Role,
		HardSkills:    p.HardSkills,
		SoftSkills:    p.SoftSkills,
		JoinTimestamp: p.JoinTimestamp,
	}
}

func ResponseFromTeam(team *Team) TeamResponse {
	leader := ResponseFromProfile(&team.Leader)
	players := make([]ProfileResponse, 0)

	for _, player := range team.Players {
		players = append(players, ResponseFromProfile(&player))
	}

	return TeamResponse{
		ID:      team.ID,
		Leader:  leader,
		Players: players,
		Mentor:  ResponseFromProfile(&team.Mentor),
	}
}

func ResponseFromLobby(lobby *Lobby) *LobbyResponse {
	players := make([]ProfileResponse, 0)
	mentors := make([]ProfileResponse, 0)
	teams := make([]TeamResponse, 0)

	for _, player := range lobby.Players {
		players = append(players, ResponseFromProfile(&player))
	}

	for _, mentor := range lobby.Mentors {
		mentors = append(mentors, ResponseFromProfile(&mentor))
	}

	for _, team := range lobby.Teams {
		teams = append(teams, ResponseFromTeam(team))
	}

	return &LobbyResponse{
		AccessCode:    lobby.AccessCode,
		Name:          lobby.Name,
		MaxHardSkills: lobby.MaxHardSkills,
		MaxSoftSkills: lobby.MaxSoftSkills,
		Status:        lobby.Status,
		Players:       players,
		Mentors:       mentors,
		Master:        ResponseFromProfile(&lobby.Master),
		Teams:         teams,
		ChooseControl: lobby.ChooseControl,
	}
}
