-- Drop employees table
DROP INDEX IF EXISTS idx_employees_full_name;
DROP INDEX IF EXISTS idx_employees_deleted_at;
DROP INDEX IF EXISTS idx_employees_manager_id;
DROP INDEX IF EXISTS idx_employees_status;
DROP INDEX IF EXISTS idx_employees_position;
DROP INDEX IF EXISTS idx_employees_department;
DROP INDEX IF EXISTS idx_employees_last_name;
DROP INDEX IF EXISTS idx_employees_first_name;
DROP INDEX IF EXISTS idx_employees_company_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_employee_id;
DROP TABLE IF EXISTS employees;