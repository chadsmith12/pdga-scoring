package fantasy

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/chadsmith12/pdga-scoring/pkgs/pdga"
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
    RoundResults []RoundResult
    HotRounds map[int][]int64
}

type RoundResult struct {
    roundNumber int32
    playerBirdies map[int64]int
    playerEaglesBetter map[int64]int
    playerBogeys map[int64]int
    playerDoubleWorse map[int64]int
}

func NewRoundResult(roundNumber int32, birdies, eaglesBetter, bogeys, doubleWorse map[int64]int) RoundResult {
    return RoundResult {
        roundNumber: roundNumber,
        playerBirdies: birdies,
        playerEaglesBetter: eaglesBetter,
        playerBogeys: bogeys,
        playerDoubleWorse: doubleWorse,
    }
}

func (r *RoundResult) BirdiesForPlayer(playerId int64) int {
    return r.playerBirdies[playerId]
}

func (r *RoundResult) EaglesBetterForPlayer(playerId int64) int {
    return r.playerEaglesBetter[playerId]
}

func (r *RoundResult) BogeysForPlayer(playerId int64) int {
    return r.playerBogeys[playerId]
}

func (r *RoundResult) DoublesWorseForPlayer(playerId int64) int {
    return r.playerDoubleWorse[playerId]
}


// Returns the number of hot rounds a certain player had in the tournament
func (r *Results) PlayerHotRounds(playerId int64) int {
    numberHotRounds := 0
    for _, round := range r.HotRounds {
        if slices.Contains(round, playerId) {
            numberHotRounds++
        } 
    }

    return numberHotRounds
}

// Creates team from a single player.
// Pass in the division the player was in to be sure they are added to the correct division for the team slot
func SingleTeam(playerId int64, division pdga.Division) CurrentTeam {
    if division == pdga.Mpo {
        return CurrentTeam {
            MpoPlayers: []int64{playerId},
            Players: []int64{playerId},
        }
    }


    return CurrentTeam {
        FpoPlayers: []int64{playerId},
        Players: []int64{playerId},
    }
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
