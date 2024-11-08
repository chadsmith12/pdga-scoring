package scoring

type PositionMap map[int]int64

type ScoringConfig struct {
    BirdieStreak StreakConfig
    StrokesBelowPar int
    StrokesAbovePar int
    HotRound int
    PositionScore PositionMap
}

type StreakConfig struct {
    Length int
    Score int
}

type FantasyTeam struct {
    Players []int64
    HotRound int64
}

func ScoreTeam(config ScoringConfig, team FantasyTeam) int {
    return 0;
}
