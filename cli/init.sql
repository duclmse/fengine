-- region DDL

DO $$ BEGIN
    IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'var_type')
        THEN CREATE TYPE VAR_TYPE AS ENUM ('i32', 'i64', 'f32', 'f64', 'bool', 'json', 'string', 'binary');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'entity_type')
        THEN CREATE TYPE ENTITY_TYPE AS ENUM ('shape', 'template', 'thing');
    END IF;
END $$;

DO $$ BEGIN
    IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'method_type')
        THEN CREATE TYPE METHOD_TYPE AS ENUM ('service', 'subscription');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS entity (
    "id"          UUID,
    "name"        VARCHAR(255) UNIQUE,
    "type"        ENTITY_TYPE,
    "description" VARCHAR(500),
    "project_id"  UUID,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS "attribute" (
    "entity_id"    UUID,
    "name"         VARCHAR(255),
    "type"         VAR_TYPE,
    "from"         UUID,
    "value_i32"    INT4,
    "value_i64"    INT4,
    "value_f32"    FLOAT4,
    "value_f64"    FLOAT8,
    "value_bool"   BOOLEAN,
    "value_json"   JSONB,
    "value_string" TEXT,
    "value_binary" BYTEA,
    PRIMARY KEY ("entity_id", "name"),
    FOREIGN KEY ("entity_id") REFERENCES entity ("id"),
    FOREIGN KEY ("from") REFERENCES entity ("id")
);

CREATE TABLE IF NOT EXISTS "method" (
    "entity_id" UUID,
    "name"      VARCHAR(255),
    "input"     JSONB,
    "output"    VAR_TYPE,
    "from"      UUID,
    "code"      TEXT,
    PRIMARY KEY ("entity_id", "name"),
    FOREIGN KEY ("entity_id") REFERENCES entity ("id"),
    FOREIGN KEY ("from") REFERENCES entity ("id")
);

DROP TABLE "method";
DROP TABLE "attribute";
DROP TABLE "entity";
DROP TYPE "var_type";
DROP TYPE "entity_type";
DROP TYPE "method_type";
-- endregion DDL

------------------------------------------------------------------------------------------------------------------------
-- region DML
-- @formatter:off

DO $$ DECLARE
    shape_id UUID := '21d2f737-31ea-4fad-a5a9-5c2fbb3e01ab'::UUID;
    templ_id UUID := '3c62e869-d806-4b7a-a770-b07c0d435452'::UUID;
    thing_id UUID := '1d6d5123-3fb8-4ab1-956f-c6f96847471d'::UUID;
    _project_id UUID := '8039354f-397a-4284-a078-f8e8ded1c6c2'::UUID;
    str_i32 VARCHAR := '[{"name":"str","type":"string"},{"name":"i32","type":"i32"}]'::jsonb;
    s_i VARCHAR := '[{"name":"s","type":"string"},{"name":"i","type":"i32"}]'::jsonb;
BEGIN
INSERT INTO entity ("id", "name", "type", "project_id")
VALUES
    (shape_id, 'Shape', 'shape', _project_id),
    (templ_id, 'Template', 'template', _project_id),
    (thing_id, 'Thing1', 'thing', _project_id);

    INSERT INTO attribute (entity_id, "name", "type", value_i32) VALUES (shape_id, 'i', 'i32', 3200);
    INSERT INTO attribute (entity_id, "name", "type", value_i64) VALUES (shape_id, 'l', 'i64', 3200);
    INSERT INTO attribute (entity_id, "name", "type", value_f32) VALUES (shape_id, 'f', 'f32', 3200);
    INSERT INTO attribute (entity_id, "name", "type", value_f64) VALUES (shape_id, 'd', 'f64', 3200);
    INSERT INTO attribute (entity_id, "name", "type", value_bool) VALUES (shape_id, 'b', 'bool', FALSE);
    INSERT INTO attribute (entity_id, "name", "type", value_json) VALUES (shape_id, 'j', 'json', '{}'::JSONB);
    INSERT INTO attribute (entity_id, "name", "type", value_string) VALUES (shape_id, 's', 'string', 'svl');
    INSERT INTO attribute (entity_id, "name", "type", value_binary) VALUES (shape_id, 'B', 'binary', 'svl'::BYTEA);

    INSERT INTO attribute (entity_id, "name", "from", "type")
    VALUES
        (templ_id, 'i', shape_id, 'i32'),
        (templ_id, 'l', shape_id, 'i64'),
        (templ_id, 'f', shape_id, 'f32'),
        (templ_id, 'd', shape_id, 'f64'),
        (templ_id, 'b', shape_id, 'bool'),
        (templ_id, 'j', shape_id, 'json'),
        (templ_id, 's', shape_id, 'string'),
        (templ_id, 'B', shape_id, 'binary');

    INSERT INTO attribute (entity_id, "name", "from", "type", value_i32) VALUES (thing_id, 'i', shape_id, 'i32', 1);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_i64) VALUES (thing_id, 'l', shape_id, 'i64', 2);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_f32) VALUES (thing_id, 'f', shape_id, 'f32', 3);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_f64) VALUES (thing_id, 'd', shape_id, 'f64', 4);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_bool) VALUES (thing_id, 'b', shape_id, 'bool', TRUE);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_json) VALUES (thing_id, 'j', shape_id, 'json', '{}'::JSONB);
    INSERT INTO attribute (entity_id, "name", "from", "type", value_string) VALUES (thing_id, 's', shape_id, 'string', 'svl');
    INSERT INTO attribute (entity_id, "name", "from", "type", value_binary) VALUES (thing_id, 'B', shape_id, 'binary', 'b'::BYTEA);

    INSERT INTO method (entity_id, "name", "input", "output", "from", code)
    VALUES
        (shape_id, 'shape_method', str_i32, 'json', null, 'return {i32:i32+1,str:str+"!"}'),
        (templ_id, 'shape_method', str_i32, 'json', shape_id, null),
        (thing_id, 'shape_method', str_i32, 'json', shape_id, null),

        (templ_id, 'templ_method', s_i, 'json', null, 'return {s:i32+"!"+s}'),
        (thing_id, 'templ_method', s_i, 'json', templ_id, null),

        (thing_id, 'thing_method', s_i, 'json', null, 'me.test({s,i});\nTable(''a'').Select({and:[{a:{$gt:10,$lt:20}}]});\nme.i=0;\nreturn {i:i+me.i, s:s+me.s}');
END $$ LANGUAGE plpgsql;
-- endregion DML
------------------------------------------------------------------------------------------------------------------------
-- @formatter:on

SELECT
    m1.entity_id, m1.name, m1."input", m1."output", m1."from",
    CASE WHEN m1."from" IS NULL THEN m1."code" ELSE m2."code" END AS code
FROM "method" m1
         LEFT JOIN "method" m2 ON m1."from" = m2.entity_id AND m1.name = m2.name
WHERE m1.entity_id = '1d6d5123-3fb8-4ab1-956f-c6f96847471d'::UUID AND m1.name = 'shape_method'

