-- Payroll Approved table (approved by Accountant)
CREATE TABLE payroll_approved (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  draft_id UUID NOT NULL UNIQUE REFERENCES payroll_drafts(id) ON DELETE CASCADE,
  fiche_paie_number VARCHAR(50) UNIQUE NOT NULL,
  accountant_id UUID NOT NULL REFERENCES users(id),
  gl_entries JSONB NOT NULL,
  approved_at TIMESTAMPTZ DEFAULT NOW(),
  digital_signature TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_payroll_approved_draft_id ON payroll_approved(draft_id);
CREATE INDEX idx_payroll_approved_fiche_paie_number ON payroll_approved(fiche_paie_number);
CREATE INDEX idx_payroll_approved_accountant_id ON payroll_approved(accountant_id);
CREATE INDEX idx_payroll_approved_approved_at ON payroll_approved(approved_at);