CREATE TABLE
    users (
              id SERIAL PRIMARY KEY,
              name varchar(64) NOT NULL,
              surname varchar(64) NOT NULL,
              patronymic varchar(64),
              nickname varchar(64) UNIQUE NOT NULL,
              email varchar(64) UNIQUE NOT NULL,
              avatar_path varchar(256),
              password_hash text NOT NULL,
              balance DECIMAL(34, 2) DEFAULT 0.00,
              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
              is_admin BOOLEAN DEFAULT FALSE,
              is_banned BOOLEAN DEFAULT FALSE
);

CREATE TYPE org_type AS ENUM ('jur', 'phys', 'ip');

CREATE type transaction_type AS ENUM (
    'org_to_project',
    'project_to_org',
    'user_to_project',
    'project_to_user',
    'user_deposit',
    'user_withdraw',
    'org_deposit',
    'org_withdraw'
);

CREATE TABLE
    organizations (
                      id SERIAL PRIMARY KEY,
                      name varchar(128) NOT NULL,
                      owner integer REFERENCES users (id) NOT NULL,
                      avatar_path varchar(256),
                      email varchar(64) NOT NULL,
                      balance DECIMAL(34, 2) DEFAULT 0.00,
                      type org_type DEFAULT NULL,
                      org_type_id integer DEFAULT null,
                      created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                      is_banned boolean DEFAULT FALSE
);

CREATE TABLE
    physical_face_project_account (
                                      id SERIAL PRIMARY KEY,
                                      BIC BIGINT NOT NULL,
                                      checking_account BIGINT NOT NULL,
                                      correspondent_account BIGINT NOT NULL,
                                      FIO varchar(256) NOT NULL,
                                      INN BIGINT NOT NULL,
                                      pasport_series integer NOT NULL,
                                      pasport_number integer NOT NULL,
                                      pasport_givenby varchar(256) NOT NULL,
                                      registration_address varchar(256) NOT NULL,
                                      post_address varchar(256) NOT NULL,
                                      pasport_page_with_photo_path varchar(256) NOT NULL,
                                      pasport_page_with_propiska_path varchar(256) NOT NULL,
                                      svid_o_postanovke_na_uchet_phys_litsa_path varchar(256) NOT NULL
);

CREATE TABLE
    juridical_face_project_accout (
                                      id SERIAL PRIMARY KEY,
                                      acts_on_base varchar(256) NOT NULL,
                                      position varchar(256) NOT NULL,
                                      BIC BIGINT NOT NULL,
                                      checking_account BIGINT NOT NULL,
                                      correspondent_account BIGINT NOT NULL,
                                      full_organisation_name varchar(128) NOT NULL,
                                      short_organisation_name varchar(128) NOT NULL,
                                      INN BIGINT NOT NULL,
                                      OGRN BIGINT NOT NULL,
                                      KPP varchar(128) NOT NULL,
                                      jur_address varchar(256) NOT NULL,
                                      fact_address varchar(256) NOT NULL,
                                      post_address varchar(256) NOT NULL,
                                      svid_o_registratsii_jur_litsa_path varchar(256) NOT NULL,
                                      svid_o_postanovke_na_nalog_uchet_path varchar(256) NOT NULL,
                                      protocol_o_nasznachenii_litsa_path varchar(256) NOT NULL,
                                      USN_path varchar(256) NOT NULL,
                                      ustav_path varchar(256) NOT NULL
);

CREATE TABLE
    ip_project_account (
                           id SERIAL PRIMARY KEY,
                           BIC BIGINT NOT NULL,
                           ras_schot BIGINT NOT NULL,
                           kor_schot BIGINT NOT NULL,
                           FIO varchar(256) NOT NULL,
                           ip_svid_serial BIGINT NOT NULL,
                           ip_svid_number BIGINT NOT NULL,
                           ip_svid_givenby varchar(256) NOT NULL,
                           INN BIGINT NOT NULL,
                           OGRN BIGINT NOT NULL,
                           jur_address varchar(256) NOT NULL,
                           fact_address varchar(256) NOT NULL,
                           post_address varchar(256) NOT NULL,
                           svid_o_postanovke_na_nalog_uchet_path varchar(256) NOT NULL,
                           ip_pasport_photo_page_path varchar(256) NOT NULL,
                           ip_pasport_propiska_path varchar(256) NOT NULL,
                           USN_path varchar(256) NOT NULL,
                           OGRNIP_path varchar(256) NOT NULL
);

CREATE TABLE
    projects (
                 id SERIAL PRIMARY KEY,
                 name varchar(128) NOT NULL,
                 creator_id integer REFERENCES organizations (id) NOT NULL,
                 quick_peek varchar(128) NOT NULL,
                 quick_peek_picture_path varchar(256),
                 content varchar(1024) NOT NULL,
                 is_public bool DEFAULT true,
                 is_completed bool DEFAULT false,
                 current_money DECIMAL(34, 2) DEFAULT 0.00,
                 wanted_money DECIMAL(34, 2) NOT NULL,
                 duration_days integer DEFAULT 30,
                 created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                 is_banned boolean DEFAULT FALSE
);

CREATE TABLE
    tags (
             id SERIAL PRIMARY KEY,
             name varchar(128) NOT NULL,
             description varchar(256),
             vector bytea
);

CREATE TABLE
    project_tags (
                     id SERIAL PRIMARY KEY,
                     project_id integer REFERENCES projects (id) NOT NULL,
                     tag_id integer REFERENCES tags (id) NOT NULL
);

CREATE TABLE
    comments (
                 id SERIAL PRIMARY KEY,
                 body varchar(1024) NOT NULL,
                 user_id integer REFERENCES users (id) NOT NULL,
                 project_id integer REFERENCES projects (id) NOT NULL,
                 created_at timestamp
);

CREATE TABLE
    user_right_at_org (
                          id SERIAL PRIMARY KEY,
                          org_id integer REFERENCES organizations (id) NOT NULL,
                          user_id integer REFERENCES users (id) NOT NULL,
                          org_account_management bool NOT NULL,
                          money_management bool NOT NULL,
                          project_management bool NOT NULL
);

CREATE TABLE
    transactions (
                     id SERIAL PRIMARY KEY,
                     from_id integer,
                     reciever_id integer,
                     type transaction_type,
                     amount DECIMAL(34, 2) NOT NULL,
                     cum_sum_of_reciever DECIMAL(34, 2),
                     cum_sum_of_sender DECIMAL(34, 2),
                     time_at timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION trg_set_project_completed()
    RETURNS trigger AS
    $$
BEGIN
        IF NEW.current_money IS NOT NULL AND NEW.wanted_money IS NOT NULL AND NEW.current_money >= NEW.wanted_money THEN
            NEW.is_completed := TRUE;
END IF;
RETURN NEW;
END;
    $$ LANGUAGE plpgsql;

CREATE INDEX idx_projects_creator_id ON projects (creator_id);
CREATE INDEX idx_organizations_type ON organizations (type);
CREATE INDEX idx_organizations_owner ON organizations (owner);
CREATE INDEX idx_comments_project_id_created_at ON comments (project_id, created_at DESC);
CREATE INDEX idx_transactions_from_id ON transactions (from_id);
CREATE INDEX idx_transactions_reciever_id ON transactions (reciever_id);

CREATE INDEX idx_project_tags_tag_project ON project_tags (tag_id, project_id) INCLUDE (id);
CREATE INDEX idx_tags_name ON tags (name);
CREATE INDEX idx_projects_public_created_desc ON projects (is_public, created_at DESC, name DESC);
