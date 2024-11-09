CREATE EXTENSION IF NOT EXISTS postgis;

ALTER TABLE event ADD COLUMN IF NOT EXISTS coordinates geography(POINT,4326);

CREATE INDEX IF NOT EXISTS event_coordinates_index ON event USING GIST ( coordinates );
