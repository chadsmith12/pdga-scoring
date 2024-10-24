-- name: CreateTournament :one
insert into tournaments (external_id, name, number_of_rounds)
values ($1, $2, $3)
returning *;
