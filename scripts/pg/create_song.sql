-- Table: prod.song

-- DROP TABLE IF EXISTS prod.song;

CREATE TABLE IF NOT EXISTS prod.song
(
    songid uuid,
    version character varying(16) COLLATE pg_catalog."default",
    title character varying(128) COLLATE pg_catalog."default",
    subtitle character varying(128) COLLATE pg_catalog."default",
    artist character varying(128) COLLATE pg_catalog."default",
    titletranslit character varying(128) COLLATE pg_catalog."default",
    subtitletranslit character varying(128) COLLATE pg_catalog."default",
    artisttranslit character varying(128) COLLATE pg_catalog."default",
    genre character varying(128) COLLATE pg_catalog."default",
    origin character varying(128) COLLATE pg_catalog."default",
    songtype character varying(128) COLLATE pg_catalog."default",
    songcategory character varying(128) COLLATE pg_catalog."default",
    banner_path character varying(256) COLLATE pg_catalog."default",
    music_path character varying(256) COLLATE pg_catalog."default",
    file_extension character varying(8) COLLATE pg_catalog."default",
    song_dir_path character varying(256) COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.song
    OWNER to postgres;