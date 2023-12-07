-- Table: prod.song_time_signature

-- DROP TABLE IF EXISTS prod.song_time_signature;

CREATE TABLE IF NOT EXISTS prod.song_time_signature
(
    songid uuid,
    time_signature_numerator integer,
    time_signature_denominator integer
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.song_time_signature
    OWNER to postgres;