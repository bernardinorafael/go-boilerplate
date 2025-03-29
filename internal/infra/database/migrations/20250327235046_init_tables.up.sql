-- Create users table
CREATE TABLE IF NOT EXISTS "users" (
	"id" VARCHAR(255) PRIMARY KEY,
	"email" VARCHAR(255) UNIQUE,
	"username" VARCHAR(255) UNIQUE,
	"name" VARCHAR(255),
	"password" VARCHAR(255),
	"avatar_url" VARCHAR(255) NULL,
	"enabled" BOOLEAN DEFAULT true,
	"locked" BOOLEAN DEFAULT false,
	"created" TIMESTAMPTZ DEFAULT now(),
	"updated" TIMESTAMPTZ DEFAULT now()
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS "sessions" (
	"id" VARCHAR(255) PRIMARY KEY,
	"user_id" VARCHAR(255),
	"ip_address" VARCHAR(255),
	"agent" VARCHAR(255),
	"refresh_token" VARCHAR(255),
	"active" BOOLEAN DEFAULT true,
	"expires" TIMESTAMPTZ,
	"created" TIMESTAMPTZ DEFAULT now(),
	"updated" TIMESTAMPTZ DEFAULT now()
);

-- Add foreign key constraint to sessions table
ALTER TABLE "sessions"
	ADD CONSTRAINT "fk_sessions_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- Create indexes for better query performance
CREATE INDEX "idx_users_created" ON "users" ("created");
CREATE INDEX "idx_sessions_active" ON "sessions" ("active");
CREATE INDEX "idx_sessions_refresh_token" ON "sessions" ("refresh_token");
