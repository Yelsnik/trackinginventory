CREATE TABLE "inventory" (
  "id"  bigserial PRIMARY KEY,
  "item" varchar NOT NULL,
  "serial_number" varchar NOT NULL,
  "price" bigint NOT NULL,
  "owner" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "inventory" ("item");

ALTER TABLE "inventory" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
