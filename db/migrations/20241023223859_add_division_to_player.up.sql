alter table players
add column division text not null check (division in ('MPO', 'FPO'))
