CREATE TYPE "public"."status" AS ENUM ('ACTIVE', 'INACTIVE', 'DELETED');

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users"
(
    "id"          bigint            NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "first_name"  text              NOT NULL,
    "last_name"   text              NOT NULL,
    "status"      "public"."status" NOT NULL DEFAULT 'INACTIVE'::status,
    "balance"     bigint            NOT NULL DEFAULT 0,
    "created_at"  timestamptz       NOT NULL DEFAULT now(),
    "modified_at" timestamptz       NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "public"."transactions"
(
    "id"         uuid        not null DEFAULT uuid_generate_v4(),
    "user_id"    bigint      NOT NULL REFERENCES "public"."users" ("id"),
    "change"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "transactions_idx__user_id" ON "public"."transactions" USING BTREE (user_id)
