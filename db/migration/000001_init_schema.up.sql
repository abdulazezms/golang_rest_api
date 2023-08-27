CREATE TABLE "t_accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar(255) UNIQUE NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "t_entries" (
  "id" integer PRIMARY KEY,
  "t_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp
);

CREATE TABLE "t_transactions" (
  "id" integer PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL
);

CREATE INDEX ON "t_accounts" ("owner");

CREATE INDEX ON "t_entries" ("t_account_id");

CREATE INDEX ON "t_transactions" ("from_account_id");

CREATE INDEX ON "t_transactions" ("to_account_id");

CREATE INDEX ON "t_transactions" ("from_account_id", "to_account_id");

ALTER TABLE "t_entries" ADD FOREIGN KEY ("t_account_id") REFERENCES "t_accounts" ("id");

ALTER TABLE "t_transactions" ADD FOREIGN KEY ("from_account_id") REFERENCES "t_accounts" ("id");

ALTER TABLE "t_transactions" ADD FOREIGN KEY ("to_account_id") REFERENCES "t_accounts" ("id");
