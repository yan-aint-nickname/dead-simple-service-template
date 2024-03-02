-- +goose Up
-- +goose StatementBegin
create table if not exists todos (
    id bigserial not null,
    name varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    primary key(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists todos;
-- +goose StatementEnd
