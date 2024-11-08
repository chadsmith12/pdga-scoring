-- name: CreateRoundScores :batchexec
insert into scores (player_id, tournament_id, layout_id, round_number, score)
values ($1, $2, $3, $4, $5);

-- name: CreateHoleScores :batchexec
insert into hole_scores (player_id, tournament_id, layout_id, round_number, hole_number, par, score_relative_to_par)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: GetPlayersHoleScores :many
select * from hole_scores hs
where tournament_id = $1
and player_id = any(sqlc.arg(player_ids)::int8[])
and (score_relative_to_par < 0 or score_relative_to_par > 0)
order by round_number, hole_number;
