create table if not exists "products" (
	"id" varchar(255) primary key,
	"name" varchar(255) unique not null,
	"price" integer not null,
	"created" timestamptz default now (),
	"updated" timestamptz default now ()
);
