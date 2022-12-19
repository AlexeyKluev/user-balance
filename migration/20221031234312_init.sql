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