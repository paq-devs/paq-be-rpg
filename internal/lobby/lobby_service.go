package lobby

import (
	"context"
	"time"

	"github.com/paq-devs/paq-be-rpg/internal/profile"
	"github.com/patrickmn/go-cache"
)

type LobbyRepository interface {
	Save(ctx context.Context, lobby *Lobby) error
	FindByAccessCode(ctx context.Context, accessCode string) (*Lobby, error)
	Update(ctx context.Context, lobby *Lobby) error
}

type LobbyService struct {
	repo  LobbyRepository
	cache *cache.Cache
}

func NewLobbyService(repo LobbyRepository) *LobbyService {
	c := cache.New(1*time.Minute, 10*time.Minute)
	return &LobbyService{
		repo:  repo,
		cache: c,
	}
}

func (service *LobbyService) CreateLobby(ctx context.Context, master profile.Profile, name string, maxHardSkills int, maxSoftSkills int) (*LobbyResponse, error) {
	lobby := NewLobby(master, name, maxHardSkills, maxSoftSkills)

	service.repo.Save(ctx, lobby)
	return ResponseFromLobby(lobby), nil
}

func (service *LobbyService) GetLobby(ctx context.Context, accessCode string) (*LobbyResponse, error) {
	if cachedLobby, found := service.cache.Get(accessCode); found {
		return cachedLobby.(*LobbyResponse), nil
	}

	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)

	return lobbyResponse, nil
}

func (service *LobbyService) JoinLobby(ctx context.Context, accessCode string, player profile.Profile) (*LobbyResponse, error) {
	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	lobby.Join(player)
	err = service.repo.Update(ctx, lobby)
	if err != nil {
		return nil, err
	}

	lobbyResponse := ResponseFromLobby(lobby)

	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)
	return lobbyResponse, nil
}

func (service *LobbyService) StartTeamCreation(ctx context.Context, accessCode string) (*LobbyResponse, error) {
	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	err = lobby.StartTeamCreation()
	if err != nil {
		return nil, err
	}

	err = service.repo.Update(ctx, lobby)
	if err != nil {
		return nil, err
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)

	go service.createTeams(lobby)
	return lobbyResponse, nil
}

func (service *LobbyService) createTeams(lobby *Lobby) {
	err := lobby.CreateTeams()

	if err != nil {
		service.moveToWaiting(context.Background(), lobby) // rollback to waiting
		return
	}

	if lobby.Status == TeamsCreated {
		err = lobby.StartLeaderTeamSelection()

		if err != nil {
			service.moveToWaiting(context.Background(), lobby) // rollback to waiting
			return
		}
	}

	err = service.repo.Update(context.Background(), lobby)
	if err != nil {
		return
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(lobby.AccessCode, lobbyResponse, cache.DefaultExpiration)
}

func (service *LobbyService) moveToWaiting(ctx context.Context, lobby *Lobby) {
	lobby.Status = Waiting

	err := service.repo.Update(ctx, lobby)
	if err != nil {
		return
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(lobby.AccessCode, lobbyResponse, cache.DefaultExpiration)
}

func (service *LobbyService) PromoteLeader(ctx context.Context, accessCode string, player profile.Profile) (*LobbyResponse, error) {
	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	err = lobby.PromoteLeader(player)
	if err != nil {
		return nil, err
	}

	if lobby.Status == TeamsCreated {
		err = lobby.StartLeaderTeamSelection()

		if err != nil {
			return nil, err
		}
	}

	err = service.repo.Update(ctx, lobby)
	if err != nil {
		return nil, err
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)
	return lobbyResponse, nil
}

func (service *LobbyService) SelectTeam(ctx context.Context, accessCode string, leader profile.Profile, teamID int) (*LobbyResponse, error) {
	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	err = lobby.SelectTeam(leader, teamID)
	if err != nil {
		return nil, err
	}

	err = service.repo.Update(ctx, lobby)
	if err != nil {
		return nil, err
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)
	return lobbyResponse, nil
}

func (service *LobbyService) SelectPlayer(ctx context.Context, accessCode string, leader profile.Profile, playerID string) (*LobbyResponse, error) {
	lobby, err := service.repo.FindByAccessCode(ctx, accessCode)
	if err != nil {
		return nil, err
	}

	if lobby == nil {
		return nil, nil
	}

	err = lobby.SelectPlayer(leader, playerID)
	if err != nil {
		return nil, err
	}

	err = service.repo.Update(ctx, lobby)
	if err != nil {
		return nil, err
	}

	lobbyResponse := ResponseFromLobby(lobby)
	service.cache.Set(accessCode, lobbyResponse, cache.DefaultExpiration)
	return lobbyResponse, nil
}
