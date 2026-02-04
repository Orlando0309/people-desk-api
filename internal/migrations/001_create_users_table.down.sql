-- Drop users table
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_employee_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;