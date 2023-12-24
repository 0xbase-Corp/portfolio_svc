CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "username" varchar,
  "email" varchar,
  "hashed_public_key" varchar,
  "signup_date" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_tags" (
  "tag_link_id" bigint PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "portfolio_id" bigint
);

ALTER TABLE "user_tags" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
