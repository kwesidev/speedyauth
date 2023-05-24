CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    uu_id uuid DEFAULT uuid_generate_v4 (),
    username VARCHAR UNIQUE NOT NULL ,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email_address VARCHAR UNIQUE NOT NULL ,
    phone_number VARCHAR NOT NULL,
    active BOOLEAN NOT NULL,
    meta_data VARCHAR,
    two_factor_enabled BOOLEAN 
);
CREATE TABLE two_factor_requests(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    request_id VARCHAR NOT NULL UNIQUE,
    ip_address VARCHAR,
    user_agent VARCHAR,
    code VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    expiry_time TIMESTAMP NOT NULL
);

CREATE TABLE user_refresh_tokens(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    token VARCHAR NOT NULL,
    created TIMESTAMP NOT NULL,
    ip_address VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    expiry_time TIMESTAMP NOT NULL
);

CREATE TABLE roles(
    id SERIAL NOT NULL PRIMARY KEY,
    type VARCHAR NOT NULL 
);
CREATE TABLE user_roles(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    role_id INTEGER REFERENCES roles(id) NOT NULL
);

INSERT INTO roles(type) VALUES('ADMIN');
INSERT INTO roles(type) VALUES('USER');

CREATE TABLE reset_password_requests(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    code VARCHAR NOT NULL UNIQUE,
    created TIMESTAMP NOT NULL,
    expiry_time TIMESTAMP NOT NULL
);
