ALTER TABLE IF EXISTS public."user"
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS second_name,
    DROP COLUMN IF EXISTS sport_types,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS description;
