DROP TABLE IF EXISTS "verify_emails" CASCADE ;

ALTER TABLE "users" DROP COLUMN "is_verified_email";

