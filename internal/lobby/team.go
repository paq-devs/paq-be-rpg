package lobby

import (
	"github.com/paq-devs/paq-be-rpg/internal/profile"
)

type Team struct {
	ID      int
	Mentor  profile.Profile
	Leader  profile.Profile
	Players []profile.Profile
}

func NewTeam(id int, mentor profile.Profile) Team {
	if mentor.Role != profile.Mentor {
		panic("mentor is not a mentor")
	}

	return Team{
		ID:     id,
		Mentor: mentor,
	}
}
