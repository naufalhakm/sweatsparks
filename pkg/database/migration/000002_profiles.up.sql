CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE profiles (
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    age SMALLINT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    location geography(Point, 4326) NOT NULL,
    interests TEXT[] NOT NULL,
    bio TEXT,
    preferences JSONB,
    birth_date DATE NOT NULL,
    last_active TIMESTAMP,
    PRIMARY KEY (user_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);  

CREATE INDEX idx_profiles_location ON profiles USING GIST(location);
CREATE INDEX idx_profiles_age ON profiles(age);
CREATE INDEX idx_profiles_gender ON profiles(gender);