DROP INDEX IF EXISTS event_coordinates_index;

ALTER TABLE event DROP COLUMN IF EXISTS coordinates;

DROP EXTENSION IF EXISTS postgis;
