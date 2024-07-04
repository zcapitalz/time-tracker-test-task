CREATE DOMAIN ksuid AS CHAR(27);

CREATE TABLE users (
    id ksuid PRIMARY KEY,    
    passport_series_and_number TEXT NOT NULL,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    address TEXT NOT NULL
);

CREATE TABLE tasks (
    id ksuid PRIMARY KEY,
    user_id ksuid REFERENCES users(id) NOT NULL
);

CREATE TYPE task_status AS ENUM ('iddle', 'in-work', 'completed');

CREATE TABLE task_status_changes (
    id BIGSERIAL PRIMARY KEY,
    task_id ksuid REFERENCES tasks(id),
    time TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    task_status task_status NOT NULL
);
CREATE INDEX idx_task_status_changes_time_task_id ON task_status_changes (task_id, time);
