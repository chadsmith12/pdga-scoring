-- name: CreateRoundScores :batchexec
insert into scores (player_id, tournament_id, layout_id, round_number, score)
values ($1, $2, $3, $4, $5);
