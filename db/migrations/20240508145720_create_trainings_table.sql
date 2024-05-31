-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trainings(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    started timestamp NOT NULL,
    finished timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_groups(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_types(
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL
);

INSERT INTO exercise_types(name) VALUES ('CARDIO'),  ('WORKOUT'), ('GYM');

CREATE TABLE IF NOT EXISTS exercises(
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL,
    description varchar(100),
    rest interval NOT NULL,
    exercise_type_id INTEGER REFERENCES exercise_types(id) NOT NULL,
    user_id INTEGER NOT NULL,
    exercise_group_id INTEGER REFERENCES exercise_groups(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_sets(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL REFERENCES exercises(id),
    weight REAL,
    reps INTEGER,
    duration interval
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercise_sets;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS exercise_groups;
DROP TABLE IF EXISTS exercise_types;
DROP TABLE IF EXISTS trainings;
-- +goose StatementEnd
