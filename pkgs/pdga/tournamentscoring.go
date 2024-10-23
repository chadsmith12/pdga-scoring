package pdga

import (
	"cmp"
	"slices"
)

type TournamentRounds []TournamentRoundData

func (rounds TournamentRounds) FinalStandings() []Score {
    finalRound := slices.MaxFunc(rounds, func(a, b TournamentRoundData) int {
        return cmp.Compare(a.Data.LiveRoundID, b.Data.LiveRoundID)
    })
    
    return finalRound.Data.Scores
}

