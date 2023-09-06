DROP TABLE IF EXISTS experiments_history, experiments, segments;

CREATE TABLE segments (
    id SERIAL PRIMARY KEY,
    slug varchar(255) NOT NULL UNIQUE
);

CREATE TABLE experiments (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    segment_id integer NOT NULL,
    started_at timestamp(0) without time zone NOT NULL,
    expired_at timestamp(0) without time zone DEFAULT NULL,

    CONSTRAINT experiments_user_segment_unique
        UNIQUE (user_id, segment_id),

    CONSTRAINT experiments_segment_fk
        FOREIGN KEY (segment_id)
        REFERENCES segments (id)
        ON DELETE CASCADE
);

CREATE TABLE experiments_history (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    segment_slug varchar(255) NOT NULL,
    operation_type varchar(255) NOT NULL,
    updated_at timestamp(0) without time zone NOT NULL
);