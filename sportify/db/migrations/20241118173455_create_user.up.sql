CREATE TABLE IF NOT EXISTS "public".user(
    id UUID NOT NULL PRIMARY KEY,
    tg_id BIGINT,
    username TEXT UNIQUE,
    password TEXT NOT NULL,
    created_at TIME WITH TIME ZONE DEFAULT NOW(),
    updated_at TIME WITH TIME ZONE DEFAULT NOW()
);

DROP TRIGGER IF EXISTS verify_updated_at_users ON public."user";
CREATE TRIGGER verify_updated_at_users
    BEFORE UPDATE
    ON public."user"
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_now();
