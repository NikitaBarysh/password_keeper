CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    login varchar(50) NOT NULL UNIQUE,
    password varchar not null
);

CREATE TABLE data
(
    user_id INT REFERENCES users(id),
    data bytea,
    event_type Varchar
);
