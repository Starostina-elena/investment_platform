ALTER TABLE user_right_at_org 
ADD CONSTRAINT user_right_at_org_org_user_unique 
UNIQUE (org_id, user_id);
