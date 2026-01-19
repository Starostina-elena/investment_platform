ALTER TABLE projects
    ADD COLUMN payback_started_date TIMESTAMP DEFAULT NULL;

ALTER TABLE projects
    ADD COLUMN money_required_to_payback DECIMAL(34, 2) DEFAULT 0.00 NOT NULL;
