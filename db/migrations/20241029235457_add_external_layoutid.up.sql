alter table if exists layouts
add column external_id int not null

-- todo: Look into removing id on layouts and just use the external_id as the unique primary key
