BEGIN;

CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    slug VARCHAR NOT NULL,
    first VARCHAR NOT NULL,
    last VARCHAR NOT NULL,
    description VARCHAR,
    birthdate DATE, 
    profile_img VARCHAR NOT NULL,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS character (
    id SERIAL PRIMARY KEY, 
    slug VARCHAR NOT NULL,
    name VARCHAR NOT NULL, 
    description VARCHAR, 
    img_name VARCHAR,
    insert_timestamp timestamp DEFAULT now(),
    person_id INT REFERENCES person(id)
);

CREATE TABLE IF NOT EXISTS creator (
    id serial primary key,
    name VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    page_url VARCHAR NOT NULL,
    description VARCHAR,
    profile_img VARCHAR, 
    date_established DATE,
    insert_timestamp timestamp DEFAULT now()
);

--CREATE TABLE IF NOT EXISTS tag (
--    id serial primary key,
--    tag VARCHAR UNIQUE
--);

CREATE TABLE IF NOT EXISTS video (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    video_url VARCHAR NOT NULL,
    slug VARCHAR NOT NULL,
    thumbnail_name VARCHAR,
    description VARCHAR,
    upload_date DATE,
    pg_rating rating,
    search_vector tsvector,
    insert_timestamp timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS video_person_rel (
    id SERIAL PRIMARY KEY,
    video_id INT references video(id) NOT NULL,
    person_id INT references person(id) NOT NULL,
    character_id INT references character(id),
    position INT,
    img VARCHAR,
    insert_timestamp timestamp DEFAULT now(),
    CONSTRAINT unique_video_person_character UNIQUE(video_id, person_id, character_id)
);

CREATE TABLE IF NOT EXISTS video_creator_rel (
    creator_id int references creator(id),
    video_id int references video(id),
    position int,
    insert_timestamp timestamp DEFAULT now(),
    PRIMARY KEY (creator_id, video_id)
);

-- CREATE TABLE video_tag_rel (
--     id serial primary key,
--     video_id int references video(id),
--     tag_id int references tag(id),
--     UNIQUE (video_id, tag_id)
-- );

COMMIT;
