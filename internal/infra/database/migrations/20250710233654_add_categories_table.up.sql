CREATE TABLE IF NOT EXISTS "categories" (
	"id" VARCHAR(255) PRIMARY KEY,
	"name" VARCHAR(255) UNIQUE NOT NULL,
	"slug" VARCHAR(255) NOT NULL,
	"active" BOOLEAN NOT NULL,
	"created_at" TIMESTAMPTZ DEFAULT now(),
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"deleted_at" TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS "idx_categories_active" ON "categories" ("active");
CREATE INDEX IF NOT EXISTS "idx_categories_slug" ON "categories" ("slug");
CREATE INDEX IF NOT EXISTS "idx_categories_deleted_at" ON "categories" ("deleted_at");
