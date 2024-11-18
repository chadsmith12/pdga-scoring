-- name: InsertFantasyRoundScores :batchexec
insert into fantasy_round_scores (player_id, tournament_id, round_number, birdies, eagles_or_better, bogeys, double_or_worse_bogeys)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertFantasyTournamentScores :batchexec
insert into fantasy_tournament_scores (player_id, tournament_id, won_event, podium_finish, top_10_finish, hot_rounds)
values ($1, $2, $3, $4, $5, $6);
