DROP TRIGGER IF EXISTS verify_updated_at_event ON public."event";
DROP FUNCTION IF EXISTS updated_at_now;

DROP TABLE IF EXISTS "public".event;

DROP TYPE IF EXISTS sport_type_enum;
DROP TYPE IF EXISTS game_level_enum;
DROP TYPE IF EXISTS creation_type_enum;