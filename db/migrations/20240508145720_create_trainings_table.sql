-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trainings(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    started timestamp NOT NULL,
    finished timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_group(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_type(
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL
);

INSERT INTO exercise_type(name) VALUES ('CARDIO'),  ('WORKOUT'), ('GYM');

CREATE TABLE IF NOT EXISTS exercise(
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL,
    description varchar(100),
    rest interval NOT NULL,
    exercise_type_id INTEGER REFERENCES exercise_type(id) NOT NULL,
    user_id INTEGER NOT NULL,
    exercise_group_id INTEGER REFERENCES exercise_group(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_set(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL REFERENCES exercise(id),
    weight REAL,
    reps INTEGER,
    duration interval
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercise_set;
DROP TABLE IF EXISTS exercise;
DROP TABLE IF EXISTS exercise_group;
DROP TABLE IF EXISTS exercise_type;
DROP TABLE IF EXISTS trainings;
-- +goose StatementEnd
