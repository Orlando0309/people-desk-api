-- Create payroll_configurations table for storing configurable payroll parameters
CREATE TABLE IF NOT EXISTS payroll_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) NOT NULL UNIQUE,
    value VARCHAR(500) NOT NULL,
    description TEXT,
    data_type VARCHAR(20) NOT NULL DEFAULT 'string' CHECK (data_type IN ('string', 'number', 'boolean')),
    category VARCHAR(50) NOT NULL DEFAULT 'general',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by UUID,
    updated_by UUID
);

-- Create index on key for faster lookups
CREATE INDEX idx_payroll_configurations_key ON payroll_configurations(key);
CREATE INDEX idx_payroll_configurations_category ON payroll_configurations(category);
CREATE INDEX idx_payroll_configurations_is_active ON payroll_configurations(is_active);

-- Insert default payroll configuration values
INSERT INTO payroll_configurations (key, value, description, data_type, category, created_by) VALUES
('cnaps_ostie_ceiling', '1600000', 'CNAPS and OSTIE contribution ceiling (8 × minimum wage)', 'number', 'contributions', gen_random_uuid()),
('cnaps_employee_rate', '0.01', 'CNAPS employee contribution rate (1%)', 'number', 'contributions', gen_random_uuid()),
('cnaps_employer_rate', '0.13', 'CNAPS employer contribution rate (13%)', 'number', 'contributions', gen_random_uuid()),
('ostie_employee_rate', '0.01', 'OSTIE employee contribution rate (1%)', 'number', 'contributions', gen_random_uuid()),
('ostie_employer_rate', '0.05', 'OSTIE employer contribution rate (5%)', 'number', 'contributions', gen_random_uuid()),
('irsa_min_tax', '2000', 'Minimum IRSA tax amount', 'number', 'tax', gen_random_uuid()),
('minimum_wage', '200000', 'Minimum wage in MGA', 'number', 'general', gen_random_uuid()),
('standard_work_hours', '8', 'Standard work hours per day', 'number', 'attendance', gen_random_uuid()),
('overtime_threshold', '8', 'Overtime threshold in hours', 'number', 'attendance', gen_random_uuid())
ON CONFLICT (key) DO NOTHING;
