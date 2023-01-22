CREATE SEQUENCE IF NOT EXISTS account_id;
CREATE SEQUENCE IF NOT EXISTS cloud_pockets;
CREATE SEQUENCE IF NOT EXISTS transactions ;

CREATE TABLE "accounts" (
    "id" int4 NOT NULL DEFAULT nextval('account_id'::regclass),
    "balance" float8 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

-- CREATE TABLE "cloud_pockets" (
--     "id" int8 NULL,
--     "name" varchar NULL,
--     "currency" varchar Null,
--     "category" varchar NULL,
--     "balance" float8 NULL,
--     "account" varchar null,
--     PRIMARY KEY ("id")
-- );

-- CREATE TABLE transactions (
-- 	"id" int8 NULL,
-- 	"refid" varchar NULL,
-- 	"pkid" int8 NULL,
-- 	"date" timestamp NULL,
-- 	"desc" varchar NULL,
-- 	"amount" float8 NULL,
-- 	"type" varchar NULL ,
--     PRIMARY KEY ("id")
-- );
