ALTER TABLE IF EXISTS public."user"
    ADD COLUMN IF NOT EXISTS first_name TEXT,
    ADD COLUMN IF NOT EXISTS second_name TEXT,
    ADD COLUMN IF NOT EXISTS sport_types sport_type_enum[],
    ADD COLUMN IF NOT EXISTS photo_url TEXT,
    ADD COLUMN IF NOT EXISTS description TEXT;
