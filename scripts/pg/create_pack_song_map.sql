-- Table: prod.pack_song_map

-- DROP TABLE IF EXISTS prod.pack_song_map;

CREATE TABLE IF NOT EXISTS prod.pack_song_map
(
    packid uuid,
    songid uuid
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.pack_song_map
    OWNER to postgres;