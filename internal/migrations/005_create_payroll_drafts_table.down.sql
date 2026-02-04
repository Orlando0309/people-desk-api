-- Drop payroll_drafts table
DROP INDEX IF EXISTS idx_payroll_drafts_employee_period;
DROP INDEX IF EXISTS idx_payroll_drafts_deleted_at;
DROP INDEX IF EXISTS idx_payroll_drafts_created_by;
DROP INDEX IF EXISTS idx_payroll_drafts_period_end;
DROP INDEX IF EXISTS idx_payroll_drafts_period_start;
DROP INDEX IF EXISTS idx_payroll_drafts_employee_id;
DROP TABLE IF EXISTS payroll_drafts;