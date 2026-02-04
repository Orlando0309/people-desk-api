-- Payroll Drafts table (created by HR)
CREATE TABLE payroll_drafts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
  gross_salary NUMERIC(15,2) NOT NULL,
  cnaps_employee NUMERIC(15,2) NOT NULL,
  cnaps_employer NUMERIC(15,2) NOT NULL,
  ostie_employee NUMERIC(15,2) NOT NULL,
  ostie_employer NUMERIC(15,2) NOT NULL,
  irsa NUMERIC(15,2) NOT NULL,
  net_salary NUMERIC(15,2) NOT NULL,
  cnaps_base NUMERIC(15,2) NOT NULL,
  ostie_base NUMERIC(15,2) NOT NULL,
  irsa_bracket VARCHAR(50),
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  UNIQUE(employee_id, period_start, period_end)
);

-- Indexes for performance
CREATE INDEX idx_payroll_drafts_employee_id ON payroll_drafts(employee_id);
CREATE INDEX idx_payroll_drafts_period_start ON payroll_drafts(period_start);
CREATE INDEX idx_payroll_drafts_period_end ON payroll_drafts(period_end);
CREATE INDEX idx_payroll_drafts_created_by ON payroll_drafts(created_by);
CREATE INDEX idx_payroll_drafts_deleted_at ON payroll_drafts(deleted_at);
CREATE INDEX idx_payroll_drafts_employee_period ON payroll_drafts(employee_id, period_start, period_end);