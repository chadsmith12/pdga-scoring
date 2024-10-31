-- name: CreateManyLayouts :batchexec
insert into layouts (id, tournament_id, name, course_name, length, units, holes, par)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;
