package lobby

import (
	"reflect"
	"testing"

	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

func TestNewTeam(t *testing.T) {
	mentor := profile.Profile{
		Role: profile.Mentor,
	}

	team := NewTeam(1, mentor)
	if team.ID != 1 {
		t.Errorf("expected team ID to be 1, got %d", team.ID)
	}
	if !reflect.DeepEqual(team.Mentor, mentor) {
		t.Errorf("expected team mentor to be mentor, got %v", team.Mentor)
	}
}

func TestNewTeam_WithNoMentor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected NewTeam to panic")
		}
	}()

	NewTeam(1, profile.Profile{})
}
