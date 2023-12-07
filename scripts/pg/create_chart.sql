-- Table: prod.chart

-- DROP TABLE IF EXISTS prod.chart;

CREATE TABLE IF NOT EXISTS prod.chart
(
    chartid uuid,
    chartname character varying(128) COLLATE pg_catalog."default",
    stepstype character varying(128) COLLATE pg_catalog."default",
    description character varying(128) COLLATE pg_catalog."default",
    chartstyle character varying(128) COLLATE pg_catalog."default",
    difficulty character varying(128) COLLATE pg_catalog."default",
    meter integer,
    credit character varying(128) COLLATE pg_catalog."default",
    stops_count integer,
    delays_count integer,
    warps_count integer,
    scrolls_count integer,
    fakes_count integer,
    songid uuid,
    speeds_count integer,
    stream real,
    voltage real,
    air real,
    "freeze" real,
    chaos real,
    taps_count integer,
    jumps_count integer,
    holds_count integer,
    mines_count integer,
    hands_count integer,
    rolls_count integer
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS prod.chart
    OWNER to postgres;