DROP TRIGGER IF EXISTS projects_set_completed ON projects;
DROP FUNCTION IF EXISTS trg_set_project_completed();

DROP INDEX IF EXISTS idx_projects_creator_id;
DROP INDEX IF EXISTS idx_organizations_type;
DROP INDEX IF EXISTS idx_organizations_owner;
DROP INDEX IF EXISTS idx_comments_project_id_created_at;
DROP INDEX IF EXISTS idx_transactions_from_id;
DROP INDEX IF EXISTS idx_transactions_reciever_id;
DROP INDEX IF EXISTS idx_project_tags_tag_project;
DROP INDEX IF EXISTS idx_tags_name;
DROP INDEX IF EXISTS idx_projects_public_created_desc;

DROP TABLE IF EXISTS project_tags;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS user_right_at_org;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS physical_face_project_account;
DROP TABLE IF EXISTS juridical_face_project_accout;
DROP TABLE IF EXISTS ip_project_account;

DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS org_type;
