-- It`s not necessary solve this problem now
-- If it will be necessary we can do something like this:
-- 1. create new type _sport_type_enum
-- 2. change in table "events" to new type _sport_type_enum
-- 3. DROP existent sport_type_enum
-- 4. CREATE TYPE new sport_type_enum
-- 5. change in table "events" to new type sport_type_enum
-- And don`t forget problem "What do with rows with existing values of enum, but that not
-- exist in old sport_type_enum?"