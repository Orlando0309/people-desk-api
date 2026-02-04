-- Drop attendance table
DROP INDEX IF EXISTS idx_attendance_employee_date;
DROP INDEX IF EXISTS idx_attendance_deleted_at;
DROP INDEX IF EXISTS idx_attendance_status;
DROP INDEX IF EXISTS idx_attendance_date;
DROP INDEX IF EXISTS idx_attendance_employee_id;
DROP TABLE IF EXISTS attendance;