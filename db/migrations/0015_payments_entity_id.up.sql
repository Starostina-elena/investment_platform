ALTER TABLE payments ADD COLUMN entity_id INT;
UPDATE payments SET entity_id = user_id WHERE entity_id IS NULL;
ALTER TABLE payments ALTER COLUMN entity_id SET NOT NULL;
ALTER TABLE payments DROP COLUMN user_id;
CREATE INDEX IF NOT EXISTS idx_payments_entity_id ON payments(entity_id);
