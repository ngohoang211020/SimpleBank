CREATE TABLE "verify_emails" (
                                 "id" bigserial PRIMARY KEY,
                                 "username" varchar NOT NULL,
                                 "email" varchar UNIQUE NOT NULL,
                                 "secret_code" varchar NOT NULL,
                                 "is_used" bool NOT NULL DEFAULT false,
                                 "created_at" timestamptz NOT NULL DEFAULT (now()),
                                 "expired_at" timestamptz NOT NULL DEFAULT (now()+ interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "users" ADD COLUMN "is_verified_email" bool not null default false;