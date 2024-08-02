CREATE TABLE IF NOT EXISTS exercise_groups(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS trainings(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    begins timestamp NOT NULL,
    finish timestamp 
);