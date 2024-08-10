package repository

import (
	lobby_ "github.com/paq-devs/paq-be-rpg/internal/lobby"
	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

type ProfileBson struct {
	ID                string              `bson:"id"`
	Name              string              `bson:"name"`
	Avatar            string              `bson:"avatar"`
	HardSkills        []profile.HardSkill `bson:"hardSkills"`
	SoftSkills        []profile.SoftSkill `bson:"softSkills"`
	Role              profile.Role        `bson:"role"`
	JoinTimestamp     int64               `bson:"joinTimestamp"`
	SelectionPriority int                 `bson:"selectionPriority"`
}

type TeamBson struct {
	ID      int           `bson:"id"`
	Mentor  ProfileBson   `bson:"mentor"`
	Leader  ProfileBson   `bson:"leader"`
	Players []ProfileBson `bson:"players"`
}

type LobbyBson struct {
	ID            string             `bson:"_id"`
	AccessCode    string             `bson:"accessCode"`
	Master        ProfileBson        `bson:"master"`
	Name          string             `bson:"name"`
	MaxHardSkills int                `bson:"maxHardSkills"`
	MaxSoftSkills int                `bson:"maxSoftSkills"`
	Players       []ProfileBson      `bson:"players"`
	Mentors       []ProfileBson      `bson:"mentors"`
	Teams         []*TeamBson        `bson:"teams"`
	Status        lobby_.LobbyStatus `bson:"status"`
	ChooseControl *ChooseControlBson `bson:"chooseControl"`
}

type ChooseControlBson struct {
	ChoosingNow ProfileBson       `bson:"choosingNow"`
	Type        lobby_.ChooseType `bson:"type"`
}

func (l *LobbyBson) ToLobby() *lobby_.Lobby {
	lobby := &lobby_.Lobby{
		ID:            l.ID,
		AccessCode:    l.AccessCode,
		Master:        l.Master.ToProfile(),
		Name:          l.Name,
		MaxHardSkills: l.MaxHardSkills,
		MaxSoftSkills: l.MaxSoftSkills,
		Players:       make([]profile.Profile, len(l.Players)),
		Mentors:       make([]profile.Profile, len(l.Mentors)),
		Teams:         make([]*lobby_.Team, len(l.Teams)),
		Status:        l.Status,
	}

	for i, player := range l.Players {
		lobby.Players[i] = player.ToProfile()
	}

	for i, mentor := range l.Mentors {
		lobby.Mentors[i] = mentor.ToProfile()
	}

	for i, team := range l.Teams {
		lobby.Teams[i] = team.ToTeam()
	}

	if l.ChooseControl != nil {
		lobby.ChooseControl = &lobby_.ChooseControl{
			ChoosingNow: l.ChooseControl.ChoosingNow.ToProfile(),
			Type:        l.ChooseControl.Type,
		}
	}

	return lobby
}

func (p *ProfileBson) ToProfile() profile.Profile {
	return profile.Profile{
		ID:                p.ID,
		Name:              p.Name,
		Avatar:            p.Avatar,
		HardSkills:        p.HardSkills,
		SoftSkills:        p.SoftSkills,
		Role:              p.Role,
		JoinTimestamp:     p.JoinTimestamp,
		SelectionPriority: p.SelectionPriority,
	}
}

func (t *TeamBson) ToTeam() *lobby_.Team {
	team := &lobby_.Team{
		ID:      t.ID,
		Mentor:  t.Mentor.ToProfile(),
		Leader:  t.Leader.ToProfile(),
		Players: make([]profile.Profile, len(t.Players)),
	}

	for i, player := range t.Players {
		team.Players[i] = player.ToProfile()
	}

	return team
}

func NewProfileBson(p profile.Profile) ProfileBson {
	return ProfileBson{
		ID:                p.ID,
		Name:              p.Name,
		Avatar:            p.Avatar,
		HardSkills:        p.HardSkills,
		SoftSkills:        p.SoftSkills,
		Role:              p.Role,
		JoinTimestamp:     p.JoinTimestamp,
		SelectionPriority: p.SelectionPriority,
	}
}

func NewTeamBson(t *lobby_.Team) *TeamBson {
	team := &TeamBson{
		ID:      t.ID,
		Mentor:  NewProfileBson(t.Mentor),
		Leader:  NewProfileBson(t.Leader),
		Players: make([]ProfileBson, len(t.Players)),
	}

	for i, player := range t.Players {
		team.Players[i] = NewProfileBson(player)
	}

	return team
}

func NewLobbyBson(l *lobby_.Lobby) LobbyBson {
	lobby := LobbyBson{
		ID:            l.ID,
		AccessCode:    l.AccessCode,
		Master:        NewProfileBson(l.Master),
		Name:          l.Name,
		MaxHardSkills: l.MaxHardSkills,
		MaxSoftSkills: l.MaxSoftSkills,
		Players:       make([]ProfileBson, len(l.Players)),
		Mentors:       make([]ProfileBson, len(l.Mentors)),
		Teams:         make([]*TeamBson, len(l.Teams)),
		Status:        l.Status,
	}

	for i, player := range l.Players {
		lobby.Players[i] = NewProfileBson(player)
	}

	for i, mentor := range l.Mentors {
		lobby.Mentors[i] = NewProfileBson(mentor)
	}

	for i, team := range l.Teams {
		lobby.Teams[i] = NewTeamBson(team)
	}

	if l.ChooseControl != nil {
		lobby.ChooseControl = &ChooseControlBson{
			ChoosingNow: NewProfileBson(l.ChooseControl.ChoosingNow),
			Type:        l.ChooseControl.Type,
		}
	}

	return lobby
}
