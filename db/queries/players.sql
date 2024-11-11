-- name: GetPlayers :many
select * from players
order by first_name;

-- name: CreateManyPlayers :batchexec
insert into players (pdga_number, first_name, last_name, name, division, city, state_prov, country)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (pdga_number) do nothing;

-- name: GetPlayersInTournament :many
SELECT player_id, p.division FROM scores
join players p on scores.player_id = p.pdga_number
where scores.tournament_id = $1 and scores.round_number = 1;
