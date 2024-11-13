package fantasy

import (
	"encoding/json"
	"io"
	"slices"
)

type Teams []Team

func UnmarshalTeams(data []byte) (Teams, error) {
	var r Teams
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Teams) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Teams) PlayerIds()[]int64 {
    playerIds := make([]int64, 0, 10)

    for _, team := range *r {
        playerIds = append(playerIds, team.Team.Players...)
    }

    return playerIds
}

func LoadTeams(reader io.Reader) (Teams, error) {
    data, err := io.ReadAll(reader) 
    if err != nil {
        return []Team{}, err
    }

    return UnmarshalTeams(data)
}

type Team struct {
	Name string   `json:"Name"`
	Team FantasyPlayers `json:"Team"`
}

type FantasyPlayers struct {
	Players []int64 `json:"Players"`
	Bench   []int64 `json:"Bench"`
}

type CurrentTeam struct {
    Name string
    MpoPlayers []int64
    FpoPlayers []int64
    Players []int64
}

type Results struct {
    MpoWinner int64
    FpoWinner int64
    Podiums []int64
    Top10s []int64
    HotRounds map[int][]int64
    RoundBirdies []map[int64]int
    RoundEaglesBetter []map[int64]int
    RoundBogeys []map[int64]int
    RoundDoubleWorse []map[int64]int
}

func (ft Team) CreateTeam(mpoPlayers []int64, fpoPlayers []int64) CurrentTeam {
    currentTeam := CurrentTeam{
        Name: ft.Name,
        MpoPlayers: make([]int64, 0, 3),
        FpoPlayers: make([]int64, 0, 2),
        Players: make([]int64, 0, 5),
    }

    for _, pdgaNumber := range ft.Team.Players {
        if slices.Contains(mpoPlayers, pdgaNumber) {
            currentTeam.MpoPlayers = append(currentTeam.MpoPlayers, pdgaNumber)
            currentTeam.Players = append(currentTeam.Players, pdgaNumber)
        } else if slices.Contains(fpoPlayers, pdgaNumber) {
            currentTeam.FpoPlayers = append(currentTeam.FpoPlayers, pdgaNumber)
            currentTeam.Players = append(currentTeam.Players, pdgaNumber)
        }
    }

    // do we still have spots left over? check the bench
    if len(currentTeam.Players) >= 5 {
        return currentTeam
    }

    if len(currentTeam.FpoPlayers) < 2 {
        currentTeam.FpoPlayers = append(currentTeam.FpoPlayers, ft.Team.Bench[0])
        currentTeam.Players = append(currentTeam.Players, ft.Team.Bench[0])
    }
    if len(currentTeam.MpoPlayers) < 3 {
        currentTeam.MpoPlayers = append(currentTeam.MpoPlayers, ft.Team.Bench[1])
        currentTeam.Players = append(currentTeam.Players, ft.Team.Bench[1])
    }

    return currentTeam
}
func (team *Team) PlayerIds() []int64 {
    return team.Team.Players
}
