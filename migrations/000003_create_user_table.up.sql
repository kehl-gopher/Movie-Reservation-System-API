CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);

--- create a trigger function to update the update field when user or movies are updated
CREATE
OR REPLACE FUNCTION trigger_update () RETURNS TRIGGER as
$$
BEGIN
UPDATE users SET NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_user_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_update();


CREATE TRIGGER update_movie_timestamp
BEFORE UPDATE ON movies
FOR EACH ROW
EXECUTE FUNCTION trigger_update();
