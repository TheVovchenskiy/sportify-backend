DO $$
    BEGIN
        IF NOT EXISTS (SELECT * FROM pg_type WHERE typname = 'sport_type_enum') THEN
            CREATE TYPE sport_type_enum AS ENUM ('volleyball', 'football', 'basketball');
        END IF;
    END
$$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT * FROM pg_type WHERE typname = 'game_level_enum') THEN
            CREATE TYPE game_level_enum AS ENUM ('low', 'mid_minus', 'mid', 'mid_plus', 'high');
        END IF;
    END
$$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT * FROM pg_type WHERE typname = 'creation_type_enum') THEN
            CREATE TYPE creation_type_enum AS ENUM ('tg', 'site');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS "public".event
(
    id uuid NOT NULL PRIMARY KEY,
    creator_id uuid,
    subscriber_ids uuid[],
    sport_type sport_type_enum NOT NULL,
    address TEXT NOT NULL,
        CONSTRAINT max_len_address CHECK (LENGTH(address) <= 4096),
    date_start TIMESTAMP WITH TIME ZONE NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL ,
    end_time TIMESTAMP WITH TIME ZONE,
    price BIGINT NOT NULL
        CONSTRAINT not_negative_price CHECK (price >= 0),
    game_level game_level_enum NOT NULL,
    description TEXT
        CONSTRAINT max_len_description CHECK (LENGTH(description) <= 16384),
    raw_message TEXT
        CONSTRAINT max_len_raw_message CHECK (LENGTH(raw_message) <= 16384),
    capacity INTEGER,
    busy INTEGER NOT NULL DEFAULT 0,
    creation_type creation_type_enum NOT NULL,
    url_message TEXT,
    url_author TEXT,
    url_preview TEXT,
    url_photos TEXT[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);


CREATE OR REPLACE FUNCTION updated_at_now()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS verify_updated_at_event ON public."event";
CREATE TRIGGER verify_updated_at_event
    BEFORE UPDATE
    ON public."event"
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_now();
