ALTER TABLE projects
    DROP COLUMN IF EXISTS monetization_type;

DROP TYPE IF EXISTS monetization_type;
