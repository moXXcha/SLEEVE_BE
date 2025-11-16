-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "firebase_uid" character varying NOT NULL,
  "email" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "user_deleted_at" to table: "users"
CREATE INDEX "user_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "user_email" to table: "users"
CREATE UNIQUE INDEX "user_email" ON "public"."users" ("email");
-- Create index "user_firebase_uid" to table: "users"
CREATE UNIQUE INDEX "user_firebase_uid" ON "public"."users" ("firebase_uid");
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
-- Create index "users_firebase_uid_key" to table: "users"
CREATE UNIQUE INDEX "users_firebase_uid_key" ON "public"."users" ("firebase_uid");
