DROP TABLE IF EXISTS experiments, segments, users;

CREATE TABLE users (
    id integer PRIMARY KEY
);

CREATE TABLE segments (
    id   SERIAL       PRIMARY KEY,
    slug varchar(255) NOT NULL UNIQUE
);

CREATE TABLE experiments (
    id          SERIAL        PRIMARY KEY,
    user_id     integer       NOT NULL,
    segment_id  integer       NOT NULL,
    started_at  timestamp(0)  NOT NULL     DEFAULT NOW(),
    expired_at  timestamp(0)               DEFAULT NULL,

    CONSTRAINT experiments_user_segment_unique
        UNIQUE (user_id, segment_id),

    CONSTRAINT experiments_segment_fk
        FOREIGN KEY (segment_id)
        REFERENCES segments (id)
        ON DELETE CASCADE
);

-- New users trigger function
CREATE OR REPLACE FUNCTION insert_new_users() RETURNS trigger AS $$
BEGIN
    INSERT INTO users (id)
        VALUES (NEW.user_id)
        ON CONFLICT (id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Simulate the presence of users table: fill it with each insertion of experiment
CREATE OR REPLACE TRIGGER new_users
    AFTER INSERT ON experiments
    FOR EACH ROW
    EXECUTE FUNCTION insert_new_users();
