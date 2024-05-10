BEGIN TRANSACTION;
CREATE TABLE if NOT EXISTS user_count (c BIGINT);
CREATE FUNCTION user_count_increment() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN IF TG_OP = 'INSERT' THEN
UPDATE user_count
SET c = c + 1;
RETURN NEW;
ELSIF TG_OP = 'DELETE' THEN
UPDATE user_count
SET c = c - 1;
RETURN OLD;
ELSE
UPDATE user_count
SET c = 0;
RETURN NULL;
END IF;
END;
$$;
CREATE CONSTRAINT TRIGGER user_count_trigger
AFTER
INSERT
    OR DELETE ON users DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE PROCEDURE user_count_increment();
-- TRUNCATE triggers must be FOR EACH STATEMENT
CREATE TRIGGER user_count_truncate_trigger
AFTER TRUNCATE ON users FOR EACH STATEMENT EXECUTE PROCEDURE user_count_increment();
-- initialize the counter table
INSERT INTO user_count (c)
SELECT COUNT(*)
FROM users;
COMMIT;