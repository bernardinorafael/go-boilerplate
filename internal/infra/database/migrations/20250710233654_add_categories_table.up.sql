create table if not exists "categories" (
	"id" varchar(255) primary key,
	"name" varchar(255) unique not null,
	"slug" varchar(255) not null,
	"active" boolean not null,
	"created_at" timestamptz default now(),
	"updated_at" timestamptz default now(),
	"deleted_at" timestamptz null
);

create index if not exists "idx_categories_active" on "categories" ("active");
create index if not exists "idx_categories_slug" on "categories" ("slug");
create index if not exists "idx_categories_deleted_at" on "categories" ("deleted_at");
