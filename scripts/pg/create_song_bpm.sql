-- Table: prod.song_bpm

-- DROP TABLE IF EXISTS prod.song_bpm;

CREATE TABLE IF NOT EXISTS prod.song_bpm
(
    songid uuid,
    song_bpm real
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.song_bpm
    OWNER to postgres;