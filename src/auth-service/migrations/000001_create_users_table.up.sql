create table users
(
    id            serial primary key,
    username      varchar(32) not null unique,
    password_hash char(60)    not null,
    email         varchar(255) not null unique,
    created_at    timestamp   not null default now(),
    updated_at    timestamp
);