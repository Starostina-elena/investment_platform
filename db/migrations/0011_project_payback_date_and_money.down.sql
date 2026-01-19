ALTER TABLE projects
    DROP COLUMN IF EXISTS money_required_to_payback;

ALTER TABLE projects
    DROP COLUMN IF EXISTS payback_started_date;
