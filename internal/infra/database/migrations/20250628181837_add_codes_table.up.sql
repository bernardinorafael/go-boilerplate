CREATE TABLE IF NOT EXISTS "codes" (
	"id" VARCHAR(255) PRIMARY KEY,
	"user_id" VARCHAR(255) NOT NULL,
	"code" VARCHAR(34) UNIQUE NOT NULL,
	"active" BOOLEAN NOT NULL,
	"attempts" INTEGER NOT NULL,
	"used_at" TIMESTAMPTZ,
	"expires_at" TIMESTAMPTZ NOT NULL,
	"created_at" TIMESTAMPTZ DEFAULT now(),
	"updated_at" TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE "codes"
	ADD CONSTRAINT "fk_userId_codes" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE INDEX "idx_codes_userId" ON "codes" ("user_id");
CREATE INDEX "idx_codes_code" ON "codes" ("code");
CREATE INDEX "idx_codes_active" ON "codes" ("active");
