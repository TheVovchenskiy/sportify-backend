ALTER TABLE "public".event ALTER COLUMN game_level TYPE game_level_enum[] USING
    (ARRAY[game_level]::game_level_enum[]);
