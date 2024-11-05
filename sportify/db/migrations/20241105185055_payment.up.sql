DO $$
    BEGIN
        IF NOT EXISTS (SELECT * FROM pg_type WHERE typname = 'payment_status_enum') THEN
            CREATE TYPE payment_status_enum AS ENUM ('paid', 'cancelled', 'pending');
        END IF;
    END
$$;


CREATE TABLE IF NOT EXISTS "public".payment (
    id uuid primary key,
    user_id uuid NOT NULL,
    event_id uuid NOT NULL,
    confirmation_url text NOT NULL CONSTRAINT not_zero_len CHECK (LENGTH(confirmation_url) > 0),
    status payment_status_enum NOT NULL,
    amount bigint NOT NULL CONSTRAINT positive_amount CHECK (amount > 0)
);

DROP TRIGGER IF EXISTS verify_updated_at_payment ON public."payment";
CREATE TRIGGER verify_updated_at_payment
    BEFORE UPDATE
    ON public."payment"
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_now();
