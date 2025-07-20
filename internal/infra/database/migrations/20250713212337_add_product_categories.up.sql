create table if not exists "product_categories" (
	"id" varchar(255) primary key,
	"product_id" varchar(255) not null,
	"category_id" varchar(255) not null,
	"created_at" timestamptz default now()
);

alter table "product_categories"
	add constraint "fk_product_categories_product_id" foreign key ("product_id") references "products"("id") on delete cascade,
	add constraint "fk_product_categories_category_id" foreign key ("category_id") references "categories"("id") on delete cascade,
	add constraint "uk_product_categories_product_category" unique ("product_id", "category_id");

create index if not exists "idx_product_categories_product_id" on "product_categories" ("product_id");
create index if not exists "idx_product_categories_category_id" on "product_categories" ("category_id");
