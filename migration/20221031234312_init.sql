CREATE TYPE "public"."status" AS ENUM ('ACTIVE', 'INACTIVE', 'DELETED');

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users"
(
    "id"          bigint            NOT NULL DEFAULT nextval('users_id_seq'::regclass),
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

CREATE INDEX CONCURRENTLY IF NOT EXISTS "transactions_idx__user_id" ON "public"."transactions" USING BTREE (user_id);

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS reservations_id_seq;

CREATE TABLE "public"."reservations"
(
    "id"         bigint NOT NULL DEFAULT nextval('reservations_id_seq'::regclass),
    "user_id"    bigint NOT NULL REFERENCES "public"."users" ("id"),
    "order_id"   bigint NOT NULL,
    "service_id" bigint NOT NULL,
    "amount"     bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS "reservations_idx__user_id__order_id__service_id" ON "public"."reservations" USING BTREE (user_id, order_id, service_id);