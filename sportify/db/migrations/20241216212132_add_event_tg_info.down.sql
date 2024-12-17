ALTER TABLE "public".event
DROP COLUMN IF EXISTS tg_chat_id,
DROP COLUMN IF EXISTS tg_message_id;
