create table tournaments (
  id bigint primary key generated always as identity,
  external_id bigint not null,
  name text not null,
  number_of_rounds int not null
);

create table layouts (
  id bigint primary key generated always as identity,
  tournament_id bigint references tournaments (id) not null,
  name text not null
);

create table players (
  id bigint primary key generated always as identity,
  name text not null,
  first_name text not null,
  last_name text not null
);

create table scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (id) not null,
  tournament_id bigint references tournaments (id) not null,
  layout_id bigint references layouts (id) not null,
  round_number int not null,
  score int not null
);

create table hole_scores (
  id bigint primary key generated always as identity,
  player_id bigint references players (id) not null,
  tournament_id bigint references tournaments (id) not null,
  layout_id bigint references layouts (id) not null,
  round_number int not null,
  hole_number int not null,
  par int not null,
  score_relative_to_par int not null
);

create index idx_hole_scores_player_id on hole_scores using btree (player_id);
create index idx_scores_player_id on scores using btree (player_id);
create index idx_layouts_tournament_id on layouts using btree (tournament_id);
create index idx_scores_player_id on scores using btree (player_id);
create index idx_scores_tournament_id on scores using btree (tournament_id);
create index idx_scores_layout_id on scores using btree (layout_id);
create index idx_scores_round_number on scores using btree (round_number);
create index idx_hole_scores_player_id on hole_scores using btree (player_id);
create index idx_hole_scores_tournament_id on hole_scores using btree (tournament_id);
create index idx_hole_scores_layout_id on hole_scores using btree (layout_id);
create index idx_hole_scores_round_number on hole_scores using btree (round_number);

create
or replace function check_round_number () returns trigger as $$
BEGIN
    IF NEW.round_number > (SELECT number_of_rounds FROM tournaments WHERE id = NEW.tournament_id) THEN
        RAISE EXCEPTION 'Round number exceeds the number of rounds for the tournament';
    END IF;
    RETURN NEW;
END;
$$ language plpgsql;

-- check this is a valid round number
create trigger check_round_number_before_insert_scores before insert
or
update on scores for each row
execute function check_round_number ();

create trigger check_round_number_before_insert_hole_scores before insert
or
update on hole_scores for each row
execute function check_round_number ();
