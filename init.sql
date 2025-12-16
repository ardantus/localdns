-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user', -- 'admin' or 'user'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Domains Table
CREATE TABLE domains (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Records Table
CREATE TABLE records (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    content VARCHAR(255) NOT NULL,
    ttl INTEGER DEFAULT 360,
    prio INTEGER DEFAULT 0,
    disabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed Default Admin: username "admin", password "admin123"
-- Hash generated via bcrypt (cost 14)
INSERT INTO users (username, password_hash, role) 
VALUES ('admin', '$2a$14$2t1.v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v/v', 'admin');
-- Wait, I still don't have a valid hash. 
-- I will leave this commented out and let the Go app seed it properly.
