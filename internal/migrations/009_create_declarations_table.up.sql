-- Create monthly_declarations table
CREATE TABLE IF NOT EXISTS monthly_declarations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
  declaration_type VARCHAR(50) NOT NULL CHECK (declaration_type IN ('monthly', 'annual')),
  declaration_period_start DATE NOT NULL,
  declaration_period_end DATE NOT NULL,
  gross_salary NUMERIC(15,2) NOT NULL,
  net_salary NUMERIC(15,2) NOT NULL,
  deductions NUMERIC(15,2) NOT NULL DEFAULT 0,
  contributions NUMERIC(15,2) NOT NULL DEFAULT 0,
  irsa_tax NUMERIC(15,2) NOT NULL DEFAULT 0,
  cnaps_ostie NUMERIC(15,2) NOT NULL DEFAULT 0,
  status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'approved', 'rejected')),
  submitted_at TIMESTAMPTZ,
  approved_at TIMESTAMPTZ,
  accountant_id UUID REFERENCES users(id),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index on employee_id for faster lookups
CREATE INDEX idx_monthly_declarations_employee_id ON monthly_declarations(employee_id);
CREATE INDEX idx_monthly_declarations_period ON monthly_declarations(declaration_period_start, declaration_period_end);
CREATE INDEX idx_monthly_declarations_status ON monthly_declarations(status);
CREATE INDEX idx_monthly_declarations_accountant_id ON monthly_declarations(accountant_id);

-- IRSA Tax Brackets configuration table
CREATE TABLE IF NOT EXISTS irsa_tax_brackets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  min_income NUMERIC(15,2) NOT NULL,
  max_income NUMERIC(15,2),
  tax_rate NUMERIC(5,2) NOT NULL,
  min_tax NUMERIC(15,2) DEFAULT 0,
  bracket_name VARCHAR(100) NOT NULL,
  is_active BOOLEAN DEFAULT true,
  sort_order INTEGER NOT NULL DEFAULT 0,
  effective_date DATE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  created_by UUID REFERENCES users(id),
  updated_by UUID REFERENCES users(id)
);

-- Insert default IRSA brackets as per PRD
-- Use decimal tax rates (fractions) and avoid subqueries for created_by/updated_by
INSERT INTO irsa_tax_brackets (min_income, max_income, tax_rate, min_tax, bracket_name, sort_order, effective_date, created_by, updated_by) VALUES
(0, 350000, 0.00, 0, 'Tranche 1 - 0%', 1, '2026-01-01', NULL, NULL),
(350001, 400000, 0.05, 0, 'Tranche 2 - 5%', 2, '2026-01-01', NULL, NULL),
(400001, 500000, 0.10, 2500, 'Tranche 3 - 10%', 3, '2026-01-01', NULL, NULL),
(500001, 600000, 0.15, 12500, 'Tranche 4 - 15%', 4, '2026-01-01', NULL, NULL),
(600001, NULL, 0.20, 27500, 'Tranche 5 - 20%', 5, '2026-01-01', NULL, NULL);

-- Indexes for declarations
CREATE INDEX idx_monthly_declarations_type ON monthly_declarations(declaration_type);
-- Status and accountant indexes already created above; avoid duplicate index creation
CREATE INDEX idx_irsa_brackets_effective_date ON irsa_tax_brackets(effective_date);
CREATE INDEX idx_irsa_brackets_active ON irsa_tax_brackets(is_active);
