CREATE TABLE withdrawals (
    id VARCHAR(36) PRIMARY KEY,
    external_id VARCHAR(255),
    entity_id INT NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_withdrawals_entity ON withdrawals(entity_type, entity_id);
CREATE INDEX idx_withdrawals_status ON withdrawals(status);
CREATE INDEX idx_withdrawals_external_id ON withdrawals(external_id);
