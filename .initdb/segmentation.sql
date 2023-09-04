CREATE TABLE IF NOT EXISTS segments (
    id SERIAL PRIMARY KEY,
    slug varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS experiments (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    segment_id integer NOT NULL,
    started_at timestamp NOT NULL,
    expired_at timestamp,

    CONSTRAINT experiments_user_segment_unique
        UNIQUE (user_id, segment_id),

    CONSTRAINT experiments_segment_fk
        FOREIGN KEY (segment_id)
        REFERENCES segments (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS experiment_history (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    segment_slug varchar(255) NOT NULL,
    operation_type varchar(255) NOT NULL,
    updated_at timestamp NOT NULL
);