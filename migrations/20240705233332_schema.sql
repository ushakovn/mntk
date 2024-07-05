-- +goose Up
-- +goose StatementBegin
create table prognosis (
    id integer primary key autoincrement,
    score real not null,
    result text not null,
    patient_name text not null,
    patient_birth datetime not null,
    created_at datetime not null default current_timestamp
);

create table questions (
    id integer primary key autoincrement,
    prognosis_id bigint not null references prognosis(id),
    label text not null,
    answer text not null,
    score real not null,
    created_at datetime not null default current_timestamp
);

create table users (
    id integer primary key autoincrement,
    name text not null,
    password text not null,
    created_at datetime not null default current_timestamp
);

insert into users(name, password)
values ('admin', 'admin');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;

drop table questions;

drop table prognosis;
-- +goose StatementEnd
