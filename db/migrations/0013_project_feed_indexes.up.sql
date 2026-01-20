CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_projects_feed_public
ON projects (created_at DESC, id ASC)
INCLUDE (monetization_type, quick_peek, quick_peek_picture_path, current_money, wanted_money, percent, payback_started, payback_started_date, money_required_to_payback)
WHERE is_public = true AND is_banned = false AND is_completed = false;

CREATE INDEX IF NOT EXISTS idx_projects_feed_public_type
ON projects (monetization_type, created_at DESC, id ASC)
INCLUDE (quick_peek, quick_peek_picture_path, current_money, wanted_money, percent, payback_started, payback_started_date, money_required_to_payback)
WHERE is_public = true AND is_banned = false AND is_completed = false;

CREATE INDEX IF NOT EXISTS idx_projects_name_trgm
ON projects USING GIN (name gin_trgm_ops)
WHERE is_public = true AND is_banned = false AND is_completed = false;
