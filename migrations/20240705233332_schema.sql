-- +goose Up
-- +goose StatementBegin
create table prognosis (
    id integer primary key autoincrement,
    score real not null,
    result text not null,
    patient_name text default 'Анонимный пациент',
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table prognosis;
drop table questions;
-- +goose StatementEnd
