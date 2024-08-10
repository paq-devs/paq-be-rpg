package lobby

import (
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

type LobbyStatus string

const (
	Waiting          LobbyStatus = "Waiting"
	CreatingTeam     LobbyStatus = "CreatingTeam"
	TeamsCreated     LobbyStatus = "TeamsCreated"
	LeaderElection   LobbyStatus = "LeaderElection"
	LeaderTeamSelect LobbyStatus = "LeaderTeamSelect"
	PlayerSelect     LobbyStatus = "PlayerSelect"
	ReadyToStart     LobbyStatus = "ReadyToStart"
)

type Lobby struct {
	ID            string
	AccessCode    string
	Master        profile.Profile
	Name          string
	MaxHardSkills int
	MaxSoftSkills int
	Players       []profile.Profile
	Mentors       []profile.Profile
	Teams         []*Team
	Status        LobbyStatus
	ChooseControl *ChooseControl
}

type ChooseType string

const (
	PromoteLeader ChooseType = "PromoteLeader" // phase where Master promotes a Leader
	SelectTeam    ChooseType = "ChooseTeam"    // phase where Leader chooses the team
	SelectPlayer  ChooseType = "ChoosePlayer"  // phase where Leader chooses the player
)

type ChooseControl struct {
	ChoosingNow profile.Profile
	Type        ChooseType
}

const (
	LeadershipWeight    = 10
	GDPWeight           = 100
	ElectedLeaderWeight = 1000
)

func NewPromoteLeaderChooseControl(p profile.Profile) (*ChooseControl, error) {
	if p.Role != profile.Master {
		return nil, errors.New("profile is not a master")
	}

	return &ChooseControl{
		Type:        PromoteLeader,
		ChoosingNow: p,
	}, nil
}

func NewSelectTeamChooseControl(p profile.Profile) (*ChooseControl, error) {
	if p.Role != profile.Leader {
		return nil, errors.New("profile is not a leader")
	}

	return &ChooseControl{
		Type:        SelectTeam,
		ChoosingNow: p,
	}, nil
}

func NewSelectPlayerChooseControl(p profile.Profile) (*ChooseControl, error) {
	if p.Role != profile.Leader {
		return nil, errors.New("profile is not a leader")
	}

	return &ChooseControl{
		Type:        SelectPlayer,
		ChoosingNow: p,
	}, nil
}

func NewLobby(master profile.Profile, name string, maxHardSkills int, maxSoftSkills int) *Lobby {
	if master.Role != profile.Master {
		panic("master is not a master")
	}

	id := uuid.New().String()

	return &Lobby{
		ID:            id,
		AccessCode:    id[:6],
		Master:        master,
		Status:        Waiting,
		Name:          name,
		MaxHardSkills: maxHardSkills,
		MaxSoftSkills: maxSoftSkills,
		Players:       []profile.Profile{},
		Mentors:       []profile.Profile{},
	}
}

func (l *Lobby) Join(p profile.Profile) error {
	if l.Status != Waiting {
		return errors.New("invalid_status")
	}

	if p.Role == profile.Mentor {
		l.Mentors = append(l.Mentors, p)
		return nil
	}

	if len(p.HardSkills) > l.MaxHardSkills || len(p.SoftSkills) > l.MaxSoftSkills {
		return errors.New("profile has too many skills")
	}

	p.Join()
	l.Players = append(l.Players, p)
	return nil
}

func (l *Lobby) StartTeamCreation() error {
	if l.Status != Waiting {
		return errors.New("invalid_status")
	}

	if len(l.Players) < 2 {
		return errors.New("not_enough_players")
	}

	if len(l.Mentors) == 0 {
		return errors.New("not_enough_mentors")
	}

	l.Status = CreatingTeam
	return nil
}

func (l *Lobby) CreateTeams() error {
	if l.Status != CreatingTeam {
		return errors.New("invalid_status")
	}

	if !l.hasSufficienteLeaders() {
		l.Status = LeaderElection
		chooseControl, err := NewPromoteLeaderChooseControl(l.Master)

		if err != nil {
			return err
		}

		l.ChooseControl = chooseControl
		return nil
	}

	for i, mentor := range l.Mentors {
		team := NewTeam(i, mentor)
		l.Teams = append(l.Teams, &team)
	}

	l.Status = TeamsCreated
	return nil
}

func (l *Lobby) PromoteLeader(p profile.Profile) error {
	if l.Status != LeaderElection {
		return errors.New("invalid_status")
	}

	player := l.getPlayer(p.ID)

	if player.Role == profile.Leader {
		return errors.New("profile_is_already_leader")
	}

	for i, profile_ := range l.Players {
		if profile_.ID == player.ID {
			l.Players[i].Role = profile.Leader
			break
		}
	}

	if l.hasSufficienteLeaders() {
		l.Status = TeamsCreated
		l.ChooseControl = nil
	}

	return nil
}

func (l *Lobby) StartLeaderTeamSelection() error {
	if l.Status != TeamsCreated {
		return errors.New("lobby is not in TeamsCreated status")
	}

	l.DefinePriorities()
	l.Status = LeaderTeamSelect

	leader, err := l.GetNextLeader()

	if err != nil {
		return err
	}

	chooseControl, err := NewSelectTeamChooseControl(*leader)
	if err != nil {
		return err
	}

	l.ChooseControl = chooseControl
	return nil
}

func (l *Lobby) GetNextLeader() (*profile.Profile, error) {
	sort.SliceStable(l.Players, func(i, j int) bool {
		return l.Players[i].SelectionPriority < l.Players[j].SelectionPriority
	})

	if l.Status != LeaderTeamSelect {
		return nil, errors.New("lobby is not in LeaderTeamSelect or PlayerSelect status")
	}

	lastPriority := -1
	if l.ChooseControl != nil {
		lastPriority = l.ChooseControl.ChoosingNow.SelectionPriority
	}

	for _, p := range l.Players {
		if p.SelectionPriority > lastPriority {
			return &p, nil
		}
	}

	return nil, nil
}

func (l *Lobby) GetNextTeamLeaderToPick() (*profile.Profile, error) {
	sort.SliceStable(l.Teams, func(i, j int) bool {
		return l.Teams[i].Leader.SelectionPriority < l.Teams[j].Leader.SelectionPriority
	})

	if l.Status != PlayerSelect {
		return nil, errors.New("lobby is not in LeaderTeamSelect or PlayerSelect status")
	}

	lastPriority := -1
	if l.ChooseControl != nil {
		lastPriority = l.ChooseControl.ChoosingNow.SelectionPriority
	}

	for _, t := range l.Teams {
		if t.Leader.SelectionPriority > lastPriority {
			return &t.Leader, nil
		}
	}

	return nil, nil
}

func (l *Lobby) DefinePriorities() {
	sort.SliceStable(l.Players, func(i, j int) bool {
		return l.Players[i].JoinTimestamp < l.Players[j].JoinTimestamp
	})

	currentPriority := 0
	for i, p := range l.Players {
		priority := calculatePriorityWeight(p, currentPriority)

		if priority < 0 {
			continue
		}

		l.Players[i].SelectionPriority = priority
		currentPriority++
	}
}

func (l *Lobby) SelectTeam(p profile.Profile, teamID int) error {
	if l.Status != LeaderTeamSelect {
		return errors.New("lobby is not in LeaderTeamSelect status")
	}

	leader := l.getPlayer(p.ID)

	if leader == nil {
		return errors.New("profile is not in the lobby")
	}

	if leader.Role != profile.Leader {
		return errors.New("profile is not a leader")
	}

	if l.ChooseControl.ChoosingNow.ID != p.ID {
		return errors.New("it is not the turn of the leader to choose")
	}

	if teamID < 0 || teamID >= len(l.Teams) {
		return errors.New("teamID is invalid")
	}

	team := l.Teams[teamID]
	team.Leader = *leader

	nextToSelect, err := l.GetNextLeader()
	if err != nil {
		return err
	}

	l.removePlayer(p.ID)
	if nextToSelect == nil || teamID == len(l.Teams)-1 {
		l.Status = PlayerSelect
		l.ChooseControl = nil // reset control

		firstLeaderToChoose, err := l.GetNextTeamLeaderToPick()
		if err != nil {
			return err
		}

		if firstLeaderToChoose == nil {
			return errors.New("error_getting_next_team_leader")
		}

		chooseControl, err := NewSelectPlayerChooseControl(*firstLeaderToChoose)
		if err != nil {
			return err
		}

		l.ChooseControl = chooseControl
		return nil
	}

	l.ChooseControl.ChoosingNow = *nextToSelect

	return nil
}

func (l *Lobby) SelectPlayer(p profile.Profile, playerID string) error {
	if l.Status != PlayerSelect {
		return errors.New("lobby is not in PlayerSelect status")
	}

	if l.ChooseControl.ChoosingNow.ID != p.ID {
		return errors.New("it is not the turn of the leader to choose")
	}

	player := l.getPlayer(playerID)

	if player == nil {
		return errors.New("playerID is invalid")
	}

	team := l.getTeamByLeaderID(p.ID)
	if team == nil {
		return errors.New("team not found")
	}

	team.Players = append(team.Players, *player)
	l.removePlayer(playerID)

	nextToSelect, err := l.GetNextTeamLeaderToPick()
	if err != nil {
		return err
	}

	if nextToSelect == nil && len(l.Players) == 0 {
		l.ChooseControl = nil
		l.Status = ReadyToStart
		return nil
	}

	if nextToSelect == nil {
		l.ChooseControl = nil //reset control
		nextToSelect, err = l.GetNextTeamLeaderToPick()
	}

	if err != nil {
		return err
	}

	if nextToSelect == nil {
		return errors.New("error_getting_next_team_leader")
	}

	chooseControl, err := NewSelectPlayerChooseControl(*nextToSelect)
	if err != nil {
		return err
	}

	l.ChooseControl = chooseControl

	return nil
}

// less than 0 means that the player is not eligible to select
// less is more eligible
func calculatePriorityWeight(p profile.Profile, currentPriority int) int {
	switch {
	case p.HasSoftSkill(profile.Leadership) && p.HasHardSkill(profile.GDP):
		return currentPriority
	case p.HasSoftSkill(profile.Leadership):
		return currentPriority + LeadershipWeight
	case p.HasHardSkill(profile.GDP):
		return currentPriority + GDPWeight
	case p.Role == profile.Leader:
		return currentPriority + ElectedLeaderWeight
	}

	return -1
}

func (l *Lobby) hasSufficienteLeaders() bool {
	return len(l.Mentors) <= len(l.getAllLeaders())
}

func (l *Lobby) getAllLeaders() []profile.Profile {
	leaders := make([]profile.Profile, 0)
	for _, p := range l.Players {
		if p.Role == profile.Leader {
			leaders = append(leaders, p)
		}
	}
	return leaders
}

func (l *Lobby) getPlayer(id string) *profile.Profile {
	for i, p := range l.Players {
		if p.ID == id {
			return &l.Players[i]
		}
	}

	return nil
}

func (l *Lobby) getTeamByLeaderID(leaderID string) *Team {
	for _, t := range l.Teams {
		if t.Leader.ID == leaderID {
			return t
		}
	}

	return nil
}

func (l *Lobby) removePlayer(id string) {
	for i, p := range l.Players {
		if p.ID == id {
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			return
		}
	}
}
