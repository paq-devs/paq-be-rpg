package lobby

import (
	"context"
	"testing"
	"time"

	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

type LobbyRepositoryMock struct {
	Memory map[string]*Lobby
}

func NewLobbyRepositoryMock() *LobbyRepositoryMock {
	return &LobbyRepositoryMock{
		Memory: make(map[string]*Lobby),
	}
}

func (r *LobbyRepositoryMock) FindByAccessCode(ctx context.Context, accessCode string) (*Lobby, error) {
	lobby, ok := r.Memory[accessCode]
	if !ok {
		return nil, nil
	}

	return lobby, nil
}

func (r *LobbyRepositoryMock) Update(ctx context.Context, lobby *Lobby) error {
	r.Memory[lobby.AccessCode] = lobby
	return nil
}

func (r *LobbyRepositoryMock) Save(ctx context.Context, lobby *Lobby) error {
	r.Memory[lobby.AccessCode] = lobby
	return nil
}

func TestCreateLobby(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if lobby.Name != "Test" {
		t.Error("lobby name is not Test")
	}

	if lobby.Master.Name != "Master" {
		t.Error("lobby master is not Master")
	}
}

func TestJoinLobby(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	player := profile.Profile{
		Name: "Player",
		Role: profile.Player,
	}

	lobby, err = service.JoinLobby(context.Background(), lobby.AccessCode, player)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if len(lobby.Players) != 1 {
		t.Error("lobby players is not 1")
	}

	if lobby.Players[0].Name != "Player" {
		t.Error("lobby player is not Player")
	}
}

func TestJoinLobbyWithMentor(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	player := profile.Profile{
		Name: "Player",
		Role: profile.Mentor,
	}

	lobby, err = service.JoinLobby(context.Background(), lobby.AccessCode, player)

	if err != nil {
		t.Error(err)
		return
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if len(lobby.Mentors) != 1 {
		t.Error("lobby mentors is not 1")
	}

	if lobby.Mentors[0].Name != "Player" {
		t.Error("lobby mentor is not Player")
	}
}

func TestStartTeamCreationService(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if lobby.AccessCode == "" {
		t.Error("lobby access code is empty")
	}

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile2)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, mentorProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, leaderProfile)

	_, _ = service.StartTeamCreation(context.Background(), lobby.AccessCode)

	<-time.After(1 * time.Second)

	lobby, err = service.GetLobby(context.Background(), lobby.AccessCode)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != LeaderTeamSelect {
		t.Error("lobby status is not LeaderTeamSelect")
	}

	if len(lobby.Teams) != 1 {
		t.Error("lobby teams is not 1")
	}
}

func TestSelectTeamService(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if lobby.AccessCode == "" {
		t.Error("lobby access code is empty")
	}

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile2)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, mentorProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, leaderProfile)

	_, _ = service.StartTeamCreation(context.Background(), lobby.AccessCode)

	<-time.After(1 * time.Second)

	lobby, err = service.GetLobby(context.Background(), lobby.AccessCode)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != LeaderTeamSelect {
		t.Error("lobby status is not LeaderTeamSelect")
	}

	if len(lobby.Teams) != 1 {
		t.Error("lobby teams is not 1")
	}

	teamID := 0

	_, err = service.SelectTeam(context.Background(), lobby.AccessCode, leaderProfile, teamID)

	if err != nil {
		t.Error(err)
	}

	<-time.After(1 * time.Second)

	lobby, err = service.GetLobby(context.Background(), lobby.AccessCode)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != PlayerSelect {
		t.Error("lobby status is not PlayerSelect")
	}
}

func TestPlayerSelectService(t *testing.T) {
	repo := NewLobbyRepositoryMock()
	service := NewLobbyService(repo)

	master := profile.Profile{
		Name: "Master",
		Role: profile.Master,
	}

	lobby, err := service.CreateLobby(context.Background(), master, "Test", 1, 1)

	if err != nil {
		t.Error(err)
	}

	if lobby == nil {
		t.Error("lobby is nil")
		return
	}

	if lobby.AccessCode == "" {
		t.Error("lobby access code is empty")
	}

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, playerProfile2)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, mentorProfile)
	_, _ = service.JoinLobby(context.Background(), lobby.AccessCode, leaderProfile)

	_, _ = service.StartTeamCreation(context.Background(), lobby.AccessCode)

	<-time.After(1 * time.Second)

	lobby, err = service.GetLobby(context.Background(), lobby.AccessCode)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != LeaderTeamSelect {
		t.Error("lobby status is not LeaderTeamSelect")
	}

	if len(lobby.Teams) != 1 {
		t.Error("lobby teams is not 1")
	}

	teamID := 0

	_, _ = service.SelectTeam(context.Background(), lobby.AccessCode, leaderProfile, teamID)

	<-time.After(1 * time.Second)

	lobby, _ = service.GetLobby(context.Background(), lobby.AccessCode)

	if lobby.Status != PlayerSelect {
		t.Error("lobby status is not PlayerSelect")
	}

	playerID := playerProfile.ID
	lobby, err = service.SelectPlayer(context.Background(), lobby.AccessCode, leaderProfile, playerID)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != PlayerSelect { // keep selecting players
		t.Error("lobby status is not PlayerSelect")
	}

	if lobby.Teams[0].Players[0].ID != playerID {
		t.Error("lobby player id is not playerID")
	}

	playerID = playerProfile2.ID
	lobby, err = service.SelectPlayer(context.Background(), lobby.AccessCode, leaderProfile, playerID)

	if err != nil {
		t.Error(err)
	}

	if lobby.Status != ReadyToStart { // Finished selecting players
		t.Error("lobby status is not ReadyToStart")
	}
}
