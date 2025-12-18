-- ============================================================================
-- LocalDNS Database Schema
-- ============================================================================
-- This file contains the complete database schema for LocalDNS Registrar.
-- It is used to initialize a fresh database. For existing databases,
-- use migration scripts or let GORM AutoMigrate handle schema updates.
-- ============================================================================

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user', -- 'admin' or 'user'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Contact Info (used for domain WHOIS data)
    contact_name TEXT DEFAULT '',
    contact_org TEXT DEFAULT '',
    contact_email TEXT DEFAULT '',
    contact_phone TEXT DEFAULT '',
    contact_address TEXT DEFAULT '',
    contact_city TEXT DEFAULT '',
    contact_state TEXT DEFAULT '',
    contact_zip TEXT DEFAULT '',
    contact_country TEXT DEFAULT ''
);

-- Domains Table
CREATE TABLE domains (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    -- Registrant Contact Info (WHOIS data)
    registrant_name VARCHAR(255) DEFAULT '',
    registrant_org VARCHAR(255) DEFAULT '',
    registrant_email VARCHAR(255) DEFAULT '',
    registrant_phone VARCHAR(255) DEFAULT '',
    registrant_address VARCHAR(255) DEFAULT '',
    registrant_city VARCHAR(255) DEFAULT '',
    registrant_state VARCHAR(255) DEFAULT '',
    registrant_zip VARCHAR(255) DEFAULT '',
    registrant_country VARCHAR(255) DEFAULT '',
    -- Admin Contact Info (defaults to Registrant if empty)
    admin_name VARCHAR(255) DEFAULT '',
    admin_org VARCHAR(255) DEFAULT '',
    admin_email VARCHAR(255) DEFAULT '',
    admin_phone VARCHAR(255) DEFAULT '',
    admin_address VARCHAR(255) DEFAULT '',
    admin_city VARCHAR(255) DEFAULT '',
    admin_state VARCHAR(255) DEFAULT '',
    admin_zip VARCHAR(255) DEFAULT '',
    admin_country VARCHAR(255) DEFAULT '',
    -- Tech Contact Info (defaults to Registrant if empty)
    tech_name VARCHAR(255) DEFAULT '',
    tech_org VARCHAR(255) DEFAULT '',
    tech_email VARCHAR(255) DEFAULT '',
    tech_phone VARCHAR(255) DEFAULT '',
    tech_address VARCHAR(255) DEFAULT '',
    tech_city VARCHAR(255) DEFAULT '',
    tech_state VARCHAR(255) DEFAULT '',
    tech_zip VARCHAR(255) DEFAULT '',
    tech_country VARCHAR(255) DEFAULT '',
    -- Status: active, expired, suspended
    status VARCHAR(20) DEFAULT 'active'
);

-- Records Table
CREATE TABLE records (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    content VARCHAR(255) NOT NULL,
    ttl INTEGER DEFAULT 360,
    prio INTEGER DEFAULT 0,
    disabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for domain_id in records table (for faster lookups)
CREATE INDEX idx_records_domain_id ON records(domain_id);

-- RegistrarConfig Table
CREATE TABLE registrar_configs (
    id BIGSERIAL PRIMARY KEY,
    registrar_name TEXT NOT NULL,
    registrar_url TEXT DEFAULT '',
    registrar_email TEXT DEFAULT '',
    registrar_phone TEXT DEFAULT '',
    registrar_iana_id TEXT DEFAULT '9999',
    abuse_contact_email TEXT DEFAULT '',
    abuse_contact_phone TEXT DEFAULT '',
    whois_server TEXT DEFAULT '',
    name_server1 TEXT DEFAULT '',
    name_server2 TEXT DEFAULT '',
    default_ttl BIGINT DEFAULT 3600,
    default_expiry BIGINT DEFAULT 365
);

-- ============================================================================
-- Notes:
-- ============================================================================
-- 1. Default admin user and registrar config are seeded by the Go application
--    on startup (see backend/main.go). The admin credentials are:
--    username: "admin"
--    password: "admin123"
--
-- 2. GORM AutoMigrate will automatically add missing columns on application
--    startup. However, for fresh installations, this init.sql ensures all
--    tables are created with the correct schema from the start.
--
-- 3. If you have an existing database that needs migration, you can:
--    a) Restart the backend service (GORM will auto-migrate)
--    b) Run the migration script: migration_add_domain_columns.sql
-- ============================================================================
