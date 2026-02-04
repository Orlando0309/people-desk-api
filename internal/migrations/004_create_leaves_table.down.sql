-- Drop leaves table
DROP INDEX IF EXISTS idx_leaves_employee_status;
DROP INDEX IF EXISTS idx_leaves_deleted_at;
DROP INDEX IF EXISTS idx_leaves_approver_id;
DROP INDEX IF EXISTS idx_leaves_end_date;
DROP INDEX IF EXISTS idx_leaves_start_date;
DROP INDEX IF EXISTS idx_leaves_status;
DROP INDEX IF EXISTS idx_leaves_leave_type;
DROP INDEX IF EXISTS idx_leaves_employee_id;
DROP TABLE IF EXISTS leaves;