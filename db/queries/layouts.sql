-- name: CreateManyLayouts :batchexec
insert into layouts (tournament_id, name, course_name, length, units, holes, par)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: GetLayoutIdsByTournament :many
select id from layouts
where tournament_id = $1;
