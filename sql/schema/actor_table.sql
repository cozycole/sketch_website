BEGIN;
CREATE TABLE IF NOT EXISTS actor (
    id serial primary key,
    first VARCHAR NOT NULL,
    last VARCHAR NOT NULL,
    birthdate DATE NOT NULL, 
    profile_img VARCHAR NOT NULL
);
COMMIT;