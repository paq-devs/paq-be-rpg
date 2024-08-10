package profile

import (
	"time"

	"github.com/google/uuid"
)

type HardSkill string

const (
	IA          HardSkill = "IA"
	GDP         HardSkill = "GDP"
	Marketing   HardSkill = "Marketing"
	English     HardSkill = "English"
	Design      HardSkill = "Design"
	Programming HardSkill = "Programming"
)

type SoftSkill string

const (
	Communication  SoftSkill = "Communication"
	Creativity     SoftSkill = "Creativity"
	Organization   SoftSkill = "Organization"
	Empathy        SoftSkill = "Empathy"
	ProblemSolving SoftSkill = "ProblemSolving"
	Collaboration  SoftSkill = "Collaboration"
	Leadership     SoftSkill = "Leadership"
	Proactivity    SoftSkill = "Proactivity"
)

type Role string

const (
	Player Role = "Player"
	Leader Role = "Leader"
	Master Role = "Master"
	Mentor Role = "Mentor"
)

type Profile struct {
	ID                string
	Name              string
	Avatar            string
	HardSkills        []HardSkill
	SoftSkills        []SoftSkill
	Role              Role
	JoinTimestamp     int64
	SelectionPriority int // 0 is the highest priority
}

func NewPlayer(name, avatar string, hardSkills []HardSkill, softSkills []SoftSkill) Profile {
	role := Player

	if hasSoftSkill(softSkills, Leadership) || hasHardSkill(hardSkills, GDP) {
		role = Leader
	}

	return Profile{
		ID:                uuid.New().String(),
		Name:              name,
		Avatar:            avatar,
		HardSkills:        hardSkills,
		SoftSkills:        softSkills,
		Role:              role,
		SelectionPriority: -1,
	}
}

func NewMaster(name, avatar string) Profile {
	return Profile{
		ID:                uuid.New().String(),
		Name:              name,
		Avatar:            avatar,
		HardSkills:        nil,
		SoftSkills:        nil,
		Role:              Master,
		JoinTimestamp:     time.Now().Unix(),
		SelectionPriority: -1,
	}
}

func NewMentor(name, avatar string) Profile {
	return Profile{
		ID:                uuid.New().String(),
		Name:              name,
		Avatar:            avatar,
		HardSkills:        nil,
		SoftSkills:        nil,
		Role:              Mentor,
		SelectionPriority: -1,
		JoinTimestamp:     time.Now().Unix(),
	}
}

func (p *Profile) Join() {
	p.JoinTimestamp = time.Now().Unix()
}

func (p *Profile) HasHardSkill(skill HardSkill) bool {
	return hasHardSkill(p.HardSkills, skill)
}

func (p *Profile) HasSoftSkill(skill SoftSkill) bool {
	return hasSoftSkill(p.SoftSkills, skill)
}

func hasSoftSkill(skills []SoftSkill, skill SoftSkill) bool {
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	return false
}

func hasHardSkill(skills []HardSkill, skill HardSkill) bool {
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	return false
}
