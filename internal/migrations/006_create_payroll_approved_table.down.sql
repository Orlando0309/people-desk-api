-- Drop payroll_approved table
DROP INDEX IF EXISTS idx_payroll_approved_approved_at;
DROP INDEX IF EXISTS idx_payroll_approved_accountant_id;
DROP INDEX IF EXISTS idx_payroll_approved_fiche_paie_number;
DROP INDEX IF EXISTS idx_payroll_approved_draft_id;
DROP TABLE IF EXISTS payroll_approved;