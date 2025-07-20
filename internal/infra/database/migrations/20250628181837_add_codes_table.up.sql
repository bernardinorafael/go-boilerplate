create table if not exists "codes" (
	"id" varchar(255) primary key,
	"user_id" varchar(255) not null,
	"code" varchar(34) unique not null,
	"active" boolean not null,
	"attempts" integer not null,
	"used_at" timestamptz,
	"expires_at" timestamptz not null,
	"created_at" timestamptz default now(),
	"updated_at" timestamptz default now()
);

alter table "codes"
	add constraint "fk_userid_codes" foreign key ("user_id") references "users" ("id") on delete cascade;

create index "idx_codes_userid" on "codes" ("user_id");
create index "idx_codes_code" on "codes" ("code");
create index "idx_codes_active" on "codes" ("active");
