-- ========================================
-- CREATE MAIN TENANTS TABLE
-- ========================================

CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL, -- Restaurant/business name
    slug VARCHAR(100) NOT NULL UNIQUE, -- URL-friendly unique identifier (burger-palace)
 
    status VARCHAR(20) NOT NULL DEFAULT 'active', --active, inactive, suspended

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), --Record creation time
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Last update time
    deleted_at TIMESTAMPTZ -- Soft delete support
)
