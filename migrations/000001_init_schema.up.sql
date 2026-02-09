CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, 
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE balances (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL DEFAULT 0.00 CHECK (amount >= 0),
    last_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL, 
    to_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,   
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL, 
    entity_id BIGINT NOT NULL,        
    action VARCHAR(50) NOT NULL,      
    details JSONB,                    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_transactions_from ON transactions(from_user_id);
CREATE INDEX idx_transactions_to ON transactions(to_user_id);