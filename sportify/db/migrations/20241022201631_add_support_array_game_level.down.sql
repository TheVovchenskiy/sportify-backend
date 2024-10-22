DO $$
    BEGIN
        -- Check if the column is not already of type 'game_level_enum'
        IF (SELECT pg_typeof(game_level) FROM public.event LIMIT 1) <> 'game_level_enum'::regtype THEN
            ALTER TABLE "public".event ALTER COLUMN game_level TYPE game_level_enum USING
                (game_level[1]);
        END IF;
    END $$;


