ALTER TABLE payments ADD COLUMN user_id INT;
UPDATE payments SET user_id = entity_id WHERE user_id IS NULL;
ALTER TABLE payments ALTER COLUMN user_id SET NOT NULL;
ALTER TABLE payments DROP COLUMN entity_id;
DROP INDEX IF EXISTS idx_payments_entity_id;
