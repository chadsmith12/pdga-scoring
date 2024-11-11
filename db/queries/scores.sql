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

-- name: GetHotRoundsForTournament :many
WITH BestScoresPerDivisionRound AS (
    SELECT p.division, s.round_number, MIN(s.score) AS best_score
    FROM scores s
    join players p on s.player_id  = p.pdga_number 
    WHERE s.tournament_id = $1 
    GROUP BY p.division, s.round_number 
)
select p.division, s.player_id, s.score, s.round_number 
from BestScoresPerDivisionRound bsd
join scores s on bsd.best_score = s.score and bsd.round_number = s.round_number 
join players p on s.player_id = p.pdga_number 
where s.tournament_id = $1 and p.division = bsd.division
order by s.round_number;

-- name: GetTop10ByTournament :many
WITH TotalScores AS (
    SELECT s.player_id, p.division, SUM(s.score) AS total_score
    FROM scores s
    JOIN players p ON s.player_id = p.pdga_number
    WHERE s.tournament_id = $1
    GROUP BY s.player_id, p.division
),
RankedScores AS (
    SELECT player_id, division, total_score, RANK() OVER (PARTITION BY division ORDER BY total_score ASC) AS rank
    FROM TotalScores
)
SELECT rs.rank, rs.player_id, rs.division, rs.total_score
FROM RankedScores rs
where rs.rank <= 10
ORDER BY rs.division, rs.rank;
