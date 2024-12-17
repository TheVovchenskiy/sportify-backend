ALTER TABLE "public".event
ADD COLUMN IF NOT EXISTS tg_chat_id bigint,
ADD COLUMN IF NOT EXISTS tg_message_id bigint;
