-- name: GetPlayers :many
select * from players
order by first_name;
