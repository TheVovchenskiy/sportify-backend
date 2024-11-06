DROP TRIGGER IF EXISTS verify_updated_at_payment ON "public".payment;

DROP TABLE IF EXISTS "public".payment;

DROP TYPE IF EXISTS payment_status_enum;

ALTER TABLE "public".event DROP COLUMN IF EXISTS user_paid_ids;
