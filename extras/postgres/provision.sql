CREATE DATABASE db_hachi;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS ltree;




CREATE FUNCTION uuid_ltree() RETURNS uuid
AS 'SELECT md5(random()::text || clock_timestamp()::text)::uuid'
LANGUAGE SQL
IMMUTABLE;


-- Table: public.edge_registry

-- DROP TABLE IF EXISTS public.edge_registry;

CREATE TABLE IF NOT EXISTS public.edge_registry
(
    id uuid,
    label text COLLATE pg_catalog."default",
    communication_channel text COLLATE pg_catalog."default",
    path ltree[],
    update_date timestamp without time zone,
    creation_date timestamp without time zone,
    hrn text COLLATE pg_catalog."default",
    longitude text COLLATE pg_catalog."default",
    latitude text COLLATE pg_catalog."default"
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.edge_registry
    OWNER to postgres;


INSERT INTO public.edge_registry(id, label, communication_channel, path)
VALUES (uuid_ltree(), 'controller', 'main.controller', array['root.owner']::ltree[]);

/*SELECT * FROM edge_registry WHERE path ? ARRAY['*.*.location.*','*.*.transactions.large']::lquery[];*/


SELECT * FROM edge_registry
         WHERE path ? ARRAY['*.location.*']::lquery[]
         AND path ? ARRAY['*.large']::lquery[]
         OR path ? ARRAY['*.sale.nothing']::lquery[]




INSERT INTO public.edge_agents_annotations (
    id, communication_channel, path)
VALUES ('691f4a0b-d167-1428-7ba9-0b43b30640dd','agents.691f4a0bd16714287ba90b43b30640dd', array['agents', 'agents.691f4a0bd16714287ba90b43b30640dd']::ltree[]);



SELECT * FROM public.edge_agents_annotations WHERE path ? ARRAY['agents.*']::lquery[]



SELECT * FROM edge_registry
WHERE path ? ARRAY['*.location.*']::lquery[]
    AND path ? ARRAY['*.large']::lquery[]
   OR path ? ARRAY['*.sale.nothing']::lquery[]
