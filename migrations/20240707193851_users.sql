-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd
