-- Drop audit_logs table
DROP INDEX IF EXISTS idx_audit_logs_module_created;
DROP INDEX IF EXISTS idx_audit_logs_user_action;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_record_id;
DROP INDEX IF EXISTS idx_audit_logs_module;
DROP INDEX IF EXISTS idx_audit_logs_action_type;
DROP INDEX IF EXISTS idx_audit_logs_user_role;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP TABLE IF EXISTS audit_logs;