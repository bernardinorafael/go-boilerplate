create table if not exists "users" (
	"id" varchar(255) primary key,
	"email" varchar(255) unique,
	"name" varchar(255),
	"created" timestamptz default now(),
	"updated" timestamptz default now()
);
