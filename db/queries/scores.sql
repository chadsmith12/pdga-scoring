-- name: CreateRoundScores :batchexec
insert into scores (player_id, tournament_id, layout_id, round_number, score)
values ($1, $2, $3, $4, $5);

-- name: CreateHoleScores :batchexec
insert into hole_scores (player_id, tournament_id, layout_id, round_number, hole_number, par, score_relative_to_par)
values ($1, $2, $3, $4, $5, $6, $7);
