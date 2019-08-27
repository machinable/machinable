
\c machinabledb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- todo: indexes
-- todo: views

CREATE TABLE app_users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL,
    admin BOOLEAN DEFAULT false,
);

CREATE TABLE app_sessions (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES app_users(id), 
    location VARCHAR, 
    mobile BOOLEAN DEFAULT false, 
    ip VARCHAR, 
    last_accessed TIMESTAMP NOT NULL, 
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
    created TIMESTAMP NOT NULL
);

CREATE TABLE project_users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    email VARCHAR,
    username VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    read BOOLEAN DEFAULT false,
    write BOOLEAN DEFAULT false,
    role VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE project_sessions (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES project_users(id),
    project_id uuid NOT NULL REFERENCES app_projects(id),
    location VARCHAR, 
    mobile BOOLEAN DEFAULT false, 
    ip VARCHAR, 
    last_accessed TIMESTAMP NOT NULL, 
    browser VARCHAR, 
    os VARCHAR
);

CREATE TABLE project_resource_definitions (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    name VARCHAR NOT NULL,
    path_name VARCHAR NOT NULL,
    parallel_read BOOLEAN DEFAULT false,
    parallel_write BOOLEAN DEFAULT false,
    create BOOLEAN DEFAULT false,
    read BOOLEAN DEFAULT false,
    update BOOLEAN DEFAULT false,
    delete BOOLEAN DEFAULT false,
    schema JSONB,
    created TIMESTAMP NOT NULL,

    UNIQUE(project_id, path_name)
);

CREATE TABLE project_resource_objects (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    resource_path uuid NOT NULL REFERENCES project_resources(path_name),
    user_id uuid REFERENCES project_users(id),
    apikey_id uuid REFERENCES project_apikeys(id),
    created TIMESTAMP NOT NULL,
    data JSONB
);

CREATE TABLE project_logs (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id),
    endpoint_type VARCHAR,
    verb VARCHAR,
    path VARCHAR,
    status_code INT NOT NULL DEFAULT -1,
    created TIMESTAMP NOT NULL,
    response_time INT NOT NULL DEFAULT -1,
    initiator VARCHAR,
    initiator_type VARCHAR,
    initiator_id VARCHAR,
    target_id VARCHAR
);

CREATE TABLE project_apikeys (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    project_id uuid NOT NULL REFERENCES app_projects(id), 
    key_hash VARCHAR NOT NULL,
    description VARCHAR,
    read BOOLEAN DEFAULT false,
    write BOOLEAN DEFAULT false,
    role VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL
);
