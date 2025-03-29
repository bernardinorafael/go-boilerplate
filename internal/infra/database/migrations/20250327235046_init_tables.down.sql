-- Drop indexes
DROP INDEX IF EXISTS "idx_sessions_refresh_token";
DROP INDEX IF EXISTS "idx_sessions_active";
DROP INDEX IF EXISTS "idx_users_created";

-- Drop foreign key constraint
ALTER TABLE "sessions"
    DROP CONSTRAINT IF EXISTS "fk_sessions_user_id";

-- Drop tables
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "users";