CREATE SEQUENCE IF NOT EXISTS account_id;
CREATE SEQUENCE IF NOT EXISTS cloud_pockets;
REATE SEQUENCE IF NOT EXISTS transactions ;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE TABLE "cloud_pockets" (
    "id" int8 NULL,
    "name" text NULL,
    "currency" text null.
    "category" text NULL,
    "balance" float8 NULL,
    "account" text null
);

CREATE TABLE transactions (
	"id" int8 NULL,
	"refid" varchar NULL,
	"pkid" int8 NULL,
	"date" timestamp NULL,
	"desc" varchar NULL,
	"amount" float8 NULL,
	"type" varchar NULL
);
