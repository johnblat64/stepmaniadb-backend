-- Table: prod.pack

-- DROP TABLE IF EXISTS prod.pack;

CREATE TABLE IF NOT EXISTS prod.pack
(
    packid uuid,
    name character varying(128) COLLATE pg_catalog."default",
    upload_timestamp timestamp without time zone,
    pack_banner_path character varying(256) COLLATE pg_catalog."default"
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.pack
    OWNER to postgres;