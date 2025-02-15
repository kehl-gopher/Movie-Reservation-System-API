-- alter user table
ALTER TABLE users
ADD COLUMN is_activated BOOL DEFAULT false;

ALTER TABLE users
ADD COLUMN is_admin BOOL DEFAULT false;
