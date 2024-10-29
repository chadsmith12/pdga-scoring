package pdga

type FullTournamentRound []TournamentRoundResponse

type RoundPlayer struct {
    Name string
    FirstName string
    LastName string
    PdgaNumber int64 
    PlayerDivison Division
    City string
    StateProv string
    Country string
}

func (tr FullTournamentRound) Players() []RoundPlayer {
    seenPlayers := make(map[int64]RoundPlayer)
    roundPlayers := make([]RoundPlayer, 0, 100)
    for _, response := range tr {
        for _, round := range response.Data.RoundData {
            for _, roundScore := range round.Scores {
                if _, ok := seenPlayers[roundScore.PDGANum]; !ok {
                    roundPlayers = append(roundPlayers, toRoundPlayer(roundScore))
                }
            }
        }
    }

    return roundPlayers
}

func toRoundPlayer(score Score) RoundPlayer {
    return RoundPlayer{
    	Name:          score.Name,
    	FirstName:     score.FirstName,
    	LastName:      score.LastName,
    	PdgaNumber:    score.PDGANum,
    	PlayerDivison: score.Division,
    	City:          score.City,
    	StateProv:     score.StateProv,
    	Country:       score.Country,
    }
}
