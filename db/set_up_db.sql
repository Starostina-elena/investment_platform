CREATE TABLE USERS(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    nickname text UNIQUE NOT NULL,
    email text UNIQUE NOT NULL,
    avatar_path text,
    password_hash text NOT NULL,
    balance DECIMAL(34, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE,
    is_banned BOOLEAN DEFAULT FALSE
)

CREATE TYPE org_type AS ENUM ('jur', 'phys', 'ip')

CREATE type transaction_type AS ENUM (
    'org_to_project',
    'project_to_org',
    'user_to_project',
    'project_to_user',
    'user_deposit',
    'user_withdraw',
    'org_deposit',
    'org_withdraw'
)

CREATE TABLE ORGANIZATIONS(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    owner FOREIGN KEY REFERENCES USERS(id) NOT NULL,
    avatar_path text,
    email text NOT NULL, 
    balance DECIMAL(34, 2) DEFAULT 0.00,
    type org_type DEFAULT NULL,
    org_type_id integer DEFAULT null,
    created_at timestamp,
    is_banned boolean
)