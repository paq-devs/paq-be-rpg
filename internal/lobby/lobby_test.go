package lobby

import (
	"reflect"
	"testing"

	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

func TestNewLobby(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")

	lobbyName := "Test Lobby"
	maxHardSkills := 1
	maxSoftSkills := 2

	lobby := NewLobby(masterProfile, lobbyName, maxHardSkills, maxSoftSkills)

	if lobby.Master.ID != masterProfile.ID {
		t.Errorf("Expected Master to be %+v, got %+v", masterProfile, lobby.Master)
	}

	if lobby.AccessCode == "" || lobby.AccessCode == " " {
		t.Errorf("Expected AccessCode to be not empty, got empty")
	}

	if lobby.Status != Waiting {
		t.Errorf("Expected Status to be Waiting, got %v", lobby.Status)
	}
	if lobby.Name != lobbyName {
		t.Errorf("Expected Name to be %s, got %s", lobbyName, lobby.Name)
	}
	if lobby.MaxHardSkills != maxHardSkills {
		t.Errorf("Expected MaxHardSkills to be %d, got %d", maxHardSkills, lobby.MaxHardSkills)
	}
	if lobby.MaxSoftSkills != maxSoftSkills {
		t.Errorf("Expected MaxSoftSkills to be %d, got %d", maxSoftSkills, lobby.MaxSoftSkills)
	}
	if !reflect.DeepEqual(lobby.Players, []profile.Profile{}) {
		t.Errorf("Expected Players to be empty, got %+v", lobby.Players)
	}
	if !reflect.DeepEqual(lobby.Mentors, []profile.Profile{}) {
		t.Errorf("Expected Mentors to be empty, got %+v", lobby.Mentors)
	}
}

func TestNewLobby_WithNoMaster(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected NewLobby to panic")
		}
	}()

	NewLobby(profile.Profile{
		Role: profile.Mentor,
	}, "Test Lobby", 1, 2)
}

func TestJoinLobby_WithMentor(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	mentorProfile := profile.NewMentor("Mentor", "avatar")

	err := lobby.Join(mentorProfile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(lobby.Mentors, []profile.Profile{mentorProfile}) {
		t.Errorf("Expected Mentors to be %+v, got %+v", []profile.Profile{mentorProfile}, lobby.Mentors)
	}
}

func TestJoinLobby_WithMentorAndPlayers(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	mentorProfile := profile.NewMentor("Mentor", "avatar")

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)

	err := lobby.Join(mentorProfile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = lobby.Join(playerProfile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(lobby.Mentors, []profile.Profile{mentorProfile}) {
		t.Errorf("Expected Mentors to be %+v, got %+v", []profile.Profile{mentorProfile}, lobby.Mentors)
	}
	if !reflect.DeepEqual(lobby.Players, []profile.Profile{playerProfile}) {
		t.Errorf("Expected Players to be %+v, got %+v", []profile.Profile{playerProfile}, lobby.Players)
	}
}

func TestJoinLobby_WithPlayerAndLeader(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)

	err := lobby.Join(playerProfile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(lobby.Players, []profile.Profile{playerProfile}) {
		t.Errorf("Expected Players to be %+v, got %+v", []profile.Profile{playerProfile}, lobby.Players)
	}

	leaderHardSkills := []profile.HardSkill{profile.GDP}
	leaderProfile := profile.NewPlayer("Leader", "avatar", leaderHardSkills, softSkills)

	err = lobby.Join(leaderProfile)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(lobby.getAllLeaders(), []profile.Profile{leaderProfile}) {
		t.Errorf("Expected Players to be %+v, got %+v", []profile.Profile{leaderProfile}, lobby.Players)
	}
}

func TestStartTeamCreation_WhenNotEnoughPlayers(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	err := lobby.StartTeamCreation()

	if err.Error() != "not_enough_players" {
		t.Errorf("Expected error, got nil")
	}
}

func TestStartTeamCreation_WhenNoHaveMentors(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)

	err := lobby.StartTeamCreation()

	if err.Error() != "not_enough_mentors" {
		t.Errorf("Expected error, got nil")
	}
}

func TestStartTeamCreation(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)

	err := lobby.StartTeamCreation()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if lobby.Status != CreatingTeam {
		t.Errorf("Expected Status to be CreatingTeam, got %v", lobby.Status)
	}
}

func TestCreateTeams_WhenNotHaveSufficientLeaders(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")
	mentorProfile2 := profile.NewMentor("Mentor2", "avatar")

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(mentorProfile2)

	_ = lobby.StartTeamCreation()

	err := lobby.CreateTeams()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if lobby.Status != LeaderElection {
		t.Errorf("Expected Status to be LeaderElection, got %v", lobby.Status)
	}
}

func TestCreateTeams(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(leaderProfile)

	_ = lobby.StartTeamCreation()

	err := lobby.CreateTeams()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if leaderProfile.Role != profile.Leader {
		t.Errorf("Expected Role to be Leader, got %v", leaderProfile.Role)
	}

	if lobby.Status != TeamsCreated {
		t.Errorf("Expected Status to be TeamsCreated, got %v", lobby.Status)
	}

	if len(lobby.Teams) != 1 {
		t.Errorf("Expected Teams to have 1 team, got %d", len(lobby.Teams))
	}

	if lobby.Teams[0].Mentor.ID != mentorProfile.ID {
		t.Errorf("Expected Leader to be %+v, got %+v", leaderProfile, lobby.Teams[0].Leader)
	}
}

func TestPromoteLeader_WhenNotInLeaderElection(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	mentorProfile := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(leaderProfile)

	_ = lobby.StartTeamCreation()
	_ = lobby.CreateTeams()

	err := lobby.PromoteLeader(leaderProfile)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPromoteLeader_WhenProfileIsNotLeader(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)

	mentorProfile := profile.NewMentor("Mentor", "avatar")
	mentorProfile2 := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(leaderProfile)
	_ = lobby.Join(mentorProfile2)

	_ = lobby.StartTeamCreation()
	_ = lobby.CreateTeams()

	err := lobby.PromoteLeader(leaderProfile)

	if lobby.Status != LeaderElection {
		t.Errorf("Expected Status to be LeaderElection, got %v", lobby.Status)
	}

	if err.Error() != "profile_is_already_leader" {
		t.Errorf("Expected error, got nil")
	}
}

func TestPromoteLeader(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	hardSkills := []profile.HardSkill{profile.English}
	softSkills := []profile.SoftSkill{profile.Communication}
	playerProfile := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)
	playerProfile2 := profile.NewPlayer("Player", "avatar", hardSkills, softSkills)

	mentorProfile := profile.NewMentor("Mentor", "avatar")
	mentorProfile2 := profile.NewMentor("Mentor", "avatar")

	leaderSoftSkill := []profile.SoftSkill{profile.Leadership}
	leaderProfile := profile.NewPlayer("Leader", "avatar", hardSkills, leaderSoftSkill)

	_ = lobby.Join(playerProfile)
	_ = lobby.Join(playerProfile2)
	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(mentorProfile2)

	_ = lobby.Join(leaderProfile)

	_ = lobby.StartTeamCreation()
	_ = lobby.CreateTeams()

	if lobby.Status != LeaderElection {
		t.Errorf("Expected Status to be LeaderElection, got %v", lobby.Status)
	}

	if playerProfile.Role != profile.Player {
		t.Errorf("Expected Role to be Player, got %v", leaderProfile.Role)
	}

	err := lobby.PromoteLeader(playerProfile)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	for _, p := range lobby.Players {
		if p.ID == playerProfile.ID {
			if p.Role != profile.Leader {
				t.Errorf("Expected Role to be Leader, got %v", p.Role)
			}
		}
	}

	if lobby.Status != TeamsCreated {
		t.Errorf("Expected Status to be TeamsCreated, got %v", lobby.Status)
	}
}

func TestStartLeaderTeamSelection_WhenNotInTeamsCreated(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	err := lobby.StartLeaderTeamSelection()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestStartLeaderTeamSelection(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	mentorProfile := profile.NewMentor("Mentor", "avatar")
	mentorProfile2 := profile.NewMentor("Mentor", "avatar")
	mentorProfile3 := profile.NewMentor("Mentor", "avatar")

	noPriorityFirstJoinProfile := profile.NewPlayer("No Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Communication})

	firstPriorityProfile := profile.NewPlayer("First Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Leadership})
	secondPriorityProfile := profile.NewPlayer("Second Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Leadership}) // By Join Time

	thirdPriorityProfile := profile.NewPlayer("Third Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Leadership})
	fourthPriorityProfile := profile.NewPlayer("Fourth Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Leadership}) // By Join Time

	fifthPriorityProfile := profile.NewPlayer("Fifth Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Communication})
	sixthPriorityProfile := profile.NewPlayer("Sixth Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Communication}) // By Join Time

	noPriorityProfile := profile.NewPlayer("No Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Communication})

	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(mentorProfile2)
	_ = lobby.Join(mentorProfile3)

	_ = lobby.Join(noPriorityFirstJoinProfile)

	_ = lobby.Join(firstPriorityProfile)

	_ = lobby.Join(thirdPriorityProfile)
	_ = lobby.Join(fourthPriorityProfile)

	_ = lobby.Join(secondPriorityProfile)

	_ = lobby.Join(fifthPriorityProfile)
	_ = lobby.Join(sixthPriorityProfile)

	_ = lobby.Join(noPriorityProfile)

	_ = lobby.StartTeamCreation()
	_ = lobby.CreateTeams()

	err := lobby.StartLeaderTeamSelection()

	// The priority is defined by the hard skill and the join time
	// less is better
	if lobby.getPlayer(firstPriorityProfile.ID).SelectionPriority >= lobby.getPlayer(secondPriorityProfile.ID).SelectionPriority {
		t.Errorf("Expected First Priority to be greater than Second Priority, got %d and %d", lobby.getPlayer(firstPriorityProfile.ID).SelectionPriority, lobby.getPlayer(secondPriorityProfile.ID).SelectionPriority)
	}

	if lobby.getPlayer(secondPriorityProfile.ID).SelectionPriority >= lobby.getPlayer(thirdPriorityProfile.ID).SelectionPriority {
		t.Errorf("Expected Second Priority to be greater than Third Priority, got %d and %d", lobby.getPlayer(secondPriorityProfile.ID).SelectionPriority, lobby.getPlayer(thirdPriorityProfile.ID).SelectionPriority)
	}

	if lobby.getPlayer(thirdPriorityProfile.ID).SelectionPriority >= lobby.getPlayer(fourthPriorityProfile.ID).SelectionPriority {
		t.Errorf("Expected Third Priority to be greater than Fourth Priority, got %d and %d", lobby.getPlayer(thirdPriorityProfile.ID).SelectionPriority, lobby.getPlayer(fourthPriorityProfile.ID).SelectionPriority)
	}

	if lobby.getPlayer(fourthPriorityProfile.ID).SelectionPriority >= lobby.getPlayer(fifthPriorityProfile.ID).SelectionPriority {
		t.Errorf("Expected Fourth Priority to be greater than Fifth Priority, got %d and %d", lobby.getPlayer(fourthPriorityProfile.ID).SelectionPriority, lobby.getPlayer(fifthPriorityProfile.ID).SelectionPriority)
	}

	if lobby.getPlayer(fifthPriorityProfile.ID).SelectionPriority >= lobby.getPlayer(sixthPriorityProfile.ID).SelectionPriority {
		t.Errorf("Expected Fifth Priority to be greater than Sixth Priority, got %d and %d", lobby.getPlayer(fifthPriorityProfile.ID).SelectionPriority, lobby.getPlayer(sixthPriorityProfile.ID).SelectionPriority)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if lobby.Status != LeaderTeamSelect {
		t.Errorf("Expected Status to be LeaderTeamSelect, got %v", lobby.Status)
	}

	if lobby.ChooseControl.ChoosingNow.ID != firstPriorityProfile.ID {
		t.Errorf("Expected ChoosingNow to be %+v, got %+v", firstPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}

	for _, p := range lobby.Players {
		if p.ID == noPriorityFirstJoinProfile.ID || p.ID == noPriorityProfile.ID {
			if p.SelectionPriority != -1 {
				t.Errorf("Expected SelectionPriority to be -1, got %d", p.SelectionPriority)
			}
		}
	}

	_ = lobby.SelectTeam(firstPriorityProfile, 0)

	if lobby.Teams[0].Leader.ID != firstPriorityProfile.ID {
		t.Errorf("Expected Leader to be %+v, got %+v", firstPriorityProfile, lobby.Teams[0].Leader)
	}

	if lobby.ChooseControl.ChoosingNow.ID != secondPriorityProfile.ID {
		t.Errorf("Expected next leader to be %+v, got %+v", secondPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}

	_ = lobby.SelectTeam(secondPriorityProfile, 1)

	if lobby.Teams[1].Leader.ID != secondPriorityProfile.ID {
		t.Errorf("Expected Leader to be %+v, got %+v", secondPriorityProfile, lobby.Teams[1].Leader)
	}

	if lobby.ChooseControl.ChoosingNow.ID != thirdPriorityProfile.ID {
		t.Errorf("Expected next leader to be %+v, got %+v", thirdPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}

	_ = lobby.SelectTeam(thirdPriorityProfile, 2)

	if lobby.Teams[2].Leader.ID != thirdPriorityProfile.ID {
		t.Errorf("Expected Leader to be %+v, got %+v", thirdPriorityProfile, lobby.Teams[2].Leader)
	}
	// Here ends the team selection because have 3 mentors and 3 teams

	if lobby.Status != PlayerSelect {
		t.Errorf("Expected Status to be PlayerSelect, got %v", lobby.Status)
	}

	//Then return to the firstPriorityProfile to select your player
	if lobby.ChooseControl.ChoosingNow.ID != firstPriorityProfile.ID {
		t.Errorf("Expected next leader to be %+v, got %+v", firstPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}

	if lobby.ChooseControl.Type != SelectPlayer {
		t.Errorf("Expected ChooseControl Type to be SelectPlayer, got %v", lobby.ChooseControl.Type)
	}
}

func TestPlayerSelect(t *testing.T) {
	masterProfile := profile.NewMaster("Master", "avatar")
	lobby := NewLobby(masterProfile, "Test Lobby", 1, 2)

	mentorProfile := profile.NewMentor("Mentor", "avatar")
	mentorProfile2 := profile.NewMentor("Mentor", "avatar")
	mentorProfile3 := profile.NewMentor("Mentor", "avatar")

	noPriorityFirstJoinProfile := profile.NewPlayer("No Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Communication})

	firstPriorityProfile := profile.NewPlayer("First Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Leadership})
	secondPriorityProfile := profile.NewPlayer("Second Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Leadership}) // By Join Time

	thirdPriorityProfile := profile.NewPlayer("Third Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Leadership})
	fourthPriorityProfile := profile.NewPlayer("Fourth Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Leadership}) // By Join Time

	fifthPriorityProfile := profile.NewPlayer("Fifth Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Communication})
	sixthPriorityProfile := profile.NewPlayer("Sixth Priority", "avatar", []profile.HardSkill{profile.GDP}, []profile.SoftSkill{profile.Communication}) // By Join Time

	noPriorityProfile := profile.NewPlayer("No Priority", "avatar", []profile.HardSkill{profile.English}, []profile.SoftSkill{profile.Communication})

	_ = lobby.Join(mentorProfile)
	_ = lobby.Join(mentorProfile2)
	_ = lobby.Join(mentorProfile3)

	_ = lobby.Join(noPriorityFirstJoinProfile)

	_ = lobby.Join(firstPriorityProfile)

	_ = lobby.Join(thirdPriorityProfile)
	_ = lobby.Join(fourthPriorityProfile)

	_ = lobby.Join(secondPriorityProfile)

	_ = lobby.Join(fifthPriorityProfile)
	_ = lobby.Join(sixthPriorityProfile)

	_ = lobby.Join(noPriorityProfile)

	_ = lobby.StartTeamCreation()
	_ = lobby.CreateTeams()

	_ = lobby.StartLeaderTeamSelection()

	_ = lobby.SelectTeam(firstPriorityProfile, 0)
	_ = lobby.SelectTeam(secondPriorityProfile, 1)
	_ = lobby.SelectTeam(thirdPriorityProfile, 2)
	// Here ends the team selection because have 3 mentors and 3 teams

	if lobby.Status != PlayerSelect {
		t.Errorf("Expected Status to be PlayerSelect, got %v", lobby.Status)
	}

	//Then return to the firstPriorityProfile to select your player
	if lobby.ChooseControl.ChoosingNow.ID != firstPriorityProfile.ID {
		t.Errorf("Expected next leader to be %+v, got %+v", firstPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}

	if lobby.ChooseControl.Type != SelectPlayer {
		t.Errorf("Expected ChooseControl Type to be SelectPlayer, got %v", lobby.ChooseControl.Type)
	}

	err := lobby.SelectPlayer(firstPriorityProfile, noPriorityFirstJoinProfile.ID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	//Player quit from the lobby and join to a team
	if lobby.getPlayer(noPriorityFirstJoinProfile.ID) != nil {
		t.Errorf("Expected Player to be nil, got %+v", lobby.getPlayer(noPriorityFirstJoinProfile.ID))
	}

	if lobby.Teams[0].Players[0].ID != noPriorityFirstJoinProfile.ID {
		t.Errorf("Expected Player to be %+v, got %+v", noPriorityFirstJoinProfile, lobby.Teams[0].Players[0])
	}

	if lobby.ChooseControl.ChoosingNow.ID != secondPriorityProfile.ID {
		t.Errorf("Expected next leader to be %+v, got %+v", secondPriorityProfile, lobby.ChooseControl.ChoosingNow)
	}
}
