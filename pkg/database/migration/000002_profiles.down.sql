DROP INDEX IF EXISTS idx_profiles_age;
DROP INDEX IF EXISTS idx_profiles_gender;
DROP INDEX IF EXISTS idx_profiles_location;

DROP TABLE IF EXISTS profiles;

DROP EXTENSION IF EXISTS postgis;