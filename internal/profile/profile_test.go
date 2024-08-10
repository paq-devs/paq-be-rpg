package profile

import (
	"reflect"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	name := "Player"
	avatar := "avatar"
	hardSkills := []HardSkill{Programming}
	softSkills := []SoftSkill{Communication}

	profile := NewPlayer(name, avatar, hardSkills, softSkills)

	if profile.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, profile.Name)
	}
	if profile.Avatar != avatar {
		t.Errorf("Expected Avatar to be %s, got %s", avatar, profile.Avatar)
	}
	if !reflect.DeepEqual(profile.HardSkills, hardSkills) {
		t.Errorf("Expected HardSkills to be %+v, got %+v", hardSkills, profile.HardSkills)
	}
	if !reflect.DeepEqual(profile.SoftSkills, softSkills) {
		t.Errorf("Expected SoftSkills to be %+v, got %+v", softSkills, profile.SoftSkills)
	}
	if profile.Role == Leader {
		t.Errorf("Expected Role to not be Leader, got %s", profile.Role)
	}
	if profile.JoinTimestamp != 0 {
		t.Errorf("Expected JoinTimestamp to be zero, got %d", profile.JoinTimestamp)
	}
	if profile.SelectionPriority != -1 {
		t.Errorf("Expected SelectionPriority to be -1, got %d", profile.SelectionPriority)
	}
}

func TestNewPlayerWithLeaderSkill(t *testing.T) {
	name := "Player"
	avatar := "avatar"
	hardSkills := []HardSkill{Programming}
	softSkills := []SoftSkill{Leadership}

	profile := NewPlayer(name, avatar, hardSkills, softSkills)

	if profile.Role != Leader {
		t.Errorf("Expected Role to be Leader, got %s", profile.Role)
	}

	if profile.SelectionPriority != -1 {
		t.Errorf("Expected SelectionPriority to be -1, got %d", profile.SelectionPriority)
	}
}

func TestNewPlayerWithGDPSkill_ThenTurnIntoLeaderRole(t *testing.T) {
	name := "Player"
	avatar := "avatar"
	hardSkills := []HardSkill{GDP}
	softSkills := []SoftSkill{Communication}

	profile := NewPlayer(name, avatar, hardSkills, softSkills)

	if profile.Role != Leader {
		t.Errorf("Expected Role to be Leader, got %s", profile.Role)
	}

	if profile.SelectionPriority != -1 {
		t.Errorf("Expected SelectionPriority to be -1, got %d", profile.SelectionPriority)
	}
}

func TestNewMentor(t *testing.T) {
	name := "Mentor"
	avatar := "avatar"

	profile := NewMentor(name, avatar)

	if profile.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, profile.Name)
	}
	if profile.Avatar != avatar {
		t.Errorf("Expected Avatar to be %s, got %s", avatar, profile.Avatar)
	}
	if profile.Role != Mentor {
		t.Errorf("Expected Role to be Mentor, got %s", profile.Role)
	}
	if profile.JoinTimestamp != 0 {
		t.Errorf("Expected JoinTimestamp to be zero, got %d", profile.JoinTimestamp)
	}
	if profile.SelectionPriority != -1 {
		t.Errorf("Expected SelectionPriority to be -1, got %d", profile.SelectionPriority)
	}
}

func TestNewMaster(t *testing.T) {
	name := "Master"
	avatar := "avatar"

	profile := NewMaster(name, avatar)

	if profile.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, profile.Name)
	}
	if profile.Avatar != avatar {
		t.Errorf("Expected Avatar to be %s, got %s", avatar, profile.Avatar)
	}
	if profile.Role != Master {
		t.Errorf("Expected Role to be Master, got %s", profile.Role)
	}
	if profile.JoinTimestamp != 0 {
		t.Errorf("Expected JoinTimestamp to be zero, got %d", profile.JoinTimestamp)
	}
	if profile.SelectionPriority != -1 {
		t.Errorf("Expected SelectionPriority to be -1, got %d", profile.SelectionPriority)
	}
}

func TestHasHardSkill(t *testing.T) {
	hardSkills := []HardSkill{Programming, GDP}

	if !hasHardSkill(hardSkills, Programming) {
		t.Errorf("Expected to have Programming skill")
	}

	if !hasHardSkill(hardSkills, GDP) {
		t.Errorf("Expected to have GDP skill")
	}

	if hasHardSkill(hardSkills, English) {
		t.Errorf("Expected to not have Communication skill")
	}
}

func TestHasSoftSkill(t *testing.T) {
	softSkills := []SoftSkill{Communication, Leadership}

	if !hasSoftSkill(softSkills, Communication) {
		t.Errorf("Expected to have Communication skill")
	}

	if !hasSoftSkill(softSkills, Leadership) {
		t.Errorf("Expected to have Leadership skill")
	}

	if hasSoftSkill(softSkills, Creativity) {
		t.Errorf("Expected to not have Creativity skill")
	}
}

func TestProfileHasSoftSkillAndHasHardSkill(t *testing.T) {
	profile := Profile{
		HardSkills: []HardSkill{Programming, GDP},
		SoftSkills: []SoftSkill{Communication, Leadership},
	}

	if !profile.HasHardSkill(Programming) {
		t.Errorf("Expected to have Programming skill")
	}

	if !profile.HasHardSkill(GDP) {
		t.Errorf("Expected to have GDP skill")
	}

	if profile.HasHardSkill(English) {
		t.Errorf("Expected to not have Communication skill")
	}

	if !profile.HasSoftSkill(Communication) {
		t.Errorf("Expected to have Communication skill")
	}

	if !profile.HasSoftSkill(Leadership) {
		t.Errorf("Expected to have Leadership skill")
	}

	if profile.HasSoftSkill(Creativity) {
		t.Errorf("Expected to not have Creativity skill")
	}
}

func TestJoin(t *testing.T) {
	player := NewPlayer("Player", "avatar", []HardSkill{Programming}, []SoftSkill{Communication})

	player.Join()

	if player.JoinTimestamp == 0 {
		t.Errorf("Expected JoinTimestamp to be greater than zero")
	}
}
