CREATE TYPE monetization_type AS ENUM ('charity', 'custom', 'fixed_percent', 'time_percent');

ALTER TABLE projects
    ADD COLUMN monetization_type monetization_type NOT NULL DEFAULT 'charity';
