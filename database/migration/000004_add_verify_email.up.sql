CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "secret_code" text NOT NULL,
  "is_used" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "users" ADD COLUMN "is_email_verified" boolean NOT NULL DEFAULT false;