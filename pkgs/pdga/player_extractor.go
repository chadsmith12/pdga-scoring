package pdga

import "context"

func FetchPlayers(client *Client, tournamentId int, division Division) ([]Player, error) {
    scores, err := client.FetchTournamentRound(context.Background(), tournamentId, 2, division) 
    if err != nil {
        return []Player{}, err
    }

    players := make([]Player, 0, len(scores.Data.Scores))
    for _, score := range scores.Data.Scores {
        player := Player{ FirstName: score.FirstName, LastName: score.LastName, Name: score.Name, PdgaNumber: int(score.PDGANum) }
        players = append(players, player)
    }

    return players, nil
}


type Player struct {
    FirstName string
    LastName string
    Name string
    PdgaNumber int
}
