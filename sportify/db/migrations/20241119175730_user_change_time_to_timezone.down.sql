ALTER TABLE "public".user
    ALTER COLUMN created_at TYPE TIME WITH TIME ZONE USING NOW(),
    ALTER COLUMN updated_at TYPE TIME WITH TIME ZONE USING NOW();