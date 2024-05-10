BEGIN TRANSACTION;
DROP TABLE user_count;
DROP FUNCTION user_count_increment();
DROP TRIGGER user_count_trigger ON users;
DROP TRIGGER user_count_truncate_trigger ON users;
COMMIT;