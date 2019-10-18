
\c testdb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE app_tiers (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR NOT NULL UNIQUE,
  cost VARCHAR NOT NULL DEFAULT '0',
  requests INTEGER,
  projects INTEGER,
  storage INTEGER
);

INSERT INTO app_tiers (id, name, cost, requests, projects, storage) VALUES 
  ('9473a732-dd95-4b98-b776-e2d77e1966fe', 'Free', '0', 1000, 3, 256),
  ('fdabaf45-bd8f-4a2d-994e-f5bf79b2034f', 'Basic', '10', 3000, 10, 5000),
  ('bbe1450f-aaf5-497b-9f20-c2c09b64ebd8', 'Professional', '30', 10000, 25, 20000);

CREATE TABLE app_users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    tier_id uuid NOT NULL REFERENCES app_tiers(id) DEFAULT '9473a732-dd95-4b98-b776-e2d77e1966fe', 
    admin BOOLEAN DEFAULT false
);

CREATE TABLE app_sessions (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES app_users(id), 
    location VARCHAR, 
    mobile BOOLEAN DEFAULT false, 
    ip VARCHAR, 
    last_accessed TIMESTAMP NOT NULL DEFAULT NOW(), 
    browser VARCHAR, 
    os VARCHAR
);

CREATE TABLE app_projects (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES app_users(id), 
    slug VARCHAR NOT NULL UNIQUE,
    name VARCHAR,
    description VARCHAR,
    icon VARCHAR,
    user_registration BOOLEAN DEFAULT false,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE view app_project_limits AS 
  SELECT p.*, u.id as account_id, t.requests 
  FROM app_users as u 
  INNER JOIN app_projects as p ON p.user_id = u.id 
  INNER JOIN app_tiers as t ON u.tier_id = t.id;

CREATE TABLE project_users_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    email VARCHAR,
    username VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    read BOOLEAN DEFAULT false,
    write BOOLEAN DEFAULT false,
    role VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE project_sessions_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES project_users_real(id),
    project_id uuid NOT NULL REFERENCES app_projects(id),
    location VARCHAR, 
    mobile BOOLEAN DEFAULT false, 
    ip VARCHAR, 
    last_accessed TIMESTAMP NOT NULL DEFAULT NOW(), 
    browser VARCHAR, 
    os VARCHAR
);

CREATE TABLE project_logs_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    endpoint_type VARCHAR,
    verb VARCHAR,
    path VARCHAR,
    status_code INT NOT NULL DEFAULT -1,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    aligned TIMESTAMP NOT NULL,
    response_time INT NOT NULL DEFAULT -1,
    initiator VARCHAR,
    initiator_type VARCHAR,
    initiator_id VARCHAR,
    target_id VARCHAR
);

CREATE TABLE project_apikeys_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id), 
    key_hash VARCHAR NOT NULL,
    description VARCHAR,
    read BOOLEAN DEFAULT false,
    write BOOLEAN DEFAULT false,
    role VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE project_resource_definitions_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    name VARCHAR NOT NULL,
    path_name VARCHAR NOT NULL,
    parallel_read BOOLEAN DEFAULT false,
    parallel_write BOOLEAN DEFAULT false,
    "create" BOOLEAN DEFAULT false,
    "read" BOOLEAN DEFAULT false,
    "update" BOOLEAN DEFAULT false,
    "delete" BOOLEAN DEFAULT false,
    schema JSONB,
    created TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE(project_id, path_name)
);

CREATE TABLE project_resource_objects_real (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    resource_path VARCHAR NOT NULL,
    creator_type VARCHAR NOT NULL,
    creator uuid,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    data JSONB
);
CREATE INDEX project_resource_objects_idx ON project_resource_objects_real (project_id, resource_path);

CREATE TABLE project_json_real(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  project_id uuid NOT NULL REFERENCES app_projects(id),
  root_key VARCHAR NOT NULL, 
  "create" BOOLEAN DEFAULT false,
  "read" BOOLEAN DEFAULT false,
  "update" BOOLEAN DEFAULT false,
  "delete" BOOLEAN DEFAULT false,
  data jsonb, 

  UNIQUE(project_id, root_key)
);
CREATE INDEX project_json_idx ON project_json_real (project_id, root_key);


/* PARTITIONING */

CREATE OR REPLACE FUNCTION create_partition_and_insert() RETURNS trigger AS
  $BODY$
    DECLARE
      partition TEXT;
    BEGIN
      partition := TG_RELNAME || '_' || MD5(NEW.project_id::VARCHAR);
      IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname=partition) THEN
        RAISE NOTICE 'A partition has been created %',partition;
        EXECUTE 'CREATE TABLE ' || partition || ' (check (project_id = ''' || NEW.project_id || ''')) INHERITS (' || TG_RELNAME || '_real' || ');';
      END IF;
      EXECUTE 'INSERT INTO ' || partition || ' SELECT(' || TG_RELNAME || ' ' || quote_literal(NEW) || ').* RETURNING id;';
      RETURN NEW;
    END;
  $BODY$
LANGUAGE plpgsql VOLATILE
COST 100;

/* project_resource_definitions */
CREATE view project_resource_definitions as select * from project_resource_definitions_real;
ALTER view project_resource_definitions ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_resource_definitions_insert_trigger
INSTEAD OF INSERT ON project_resource_definitions
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_resource_objects */
CREATE view project_resource_objects as select * from project_resource_objects_real;
ALTER view project_resource_objects ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_resource_objects_insert_trigger
INSTEAD OF INSERT ON project_resource_objects
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_json */
CREATE view project_json as select * from project_json_real;
ALTER view project_json ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_json
INSTEAD OF INSERT ON project_json
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_apikeys */
CREATE view project_apikeys as select * from project_apikeys_real;
ALTER view project_apikeys ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_apikeys_insert_trigger
INSTEAD OF INSERT ON project_apikeys
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_logs */
CREATE view project_logs as select * from project_logs_real;
ALTER view project_logs ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_logs_insert_trigger
INSTEAD OF INSERT ON project_logs
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_sessions */
CREATE view project_sessions as select * from project_sessions_real;
ALTER view project_sessions ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_sessions_insert_trigger
INSTEAD OF INSERT ON project_sessions
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();

/* project_users */
CREATE view project_users as select * from project_users_real;
ALTER view project_users ALTER column id set DEFAULT uuid_generate_v4();
CREATE TRIGGER project_users_insert_trigger
INSTEAD OF INSERT ON project_users
FOR EACH ROW EXECUTE PROCEDURE create_partition_and_insert();
