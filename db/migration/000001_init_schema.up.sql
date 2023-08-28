CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar(255) UNIQUE NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entry" (
  "id" integer PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp
);

CREATE TABLE "transfer" (
  "id" integer PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entry" ("account_id");

CREATE INDEX ON "transfer" ("from_account_id");

CREATE INDEX ON "transfer" ("to_account_id");

CREATE INDEX ON "transfer" ("from_account_id", "to_account_id");

ALTER TABLE "entry" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("from_account_id") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("to_account_id") REFERENCES "account" ("id");
