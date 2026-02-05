-- Company settings table
CREATE TABLE company_settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_name VARCHAR(255) NOT NULL,
  company_address TEXT,
  company_nif VARCHAR(50),
  company_stat VARCHAR(50),
  cnaps_number VARCHAR(50),
  ostie_number VARCHAR(50),
  contact_email VARCHAR(255),
  contact_phone VARCHAR(50),
  logo_url TEXT,
  timezone VARCHAR(100) DEFAULT 'Indian/Antananarivo',
  currency VARCHAR(10) DEFAULT 'MGA',
  fiscal_year_start VARCHAR(5) DEFAULT '01-01',
  work_hours_per_day DECIMAL(4,2) DEFAULT 8.00,
  work_days_per_week INT DEFAULT 5,
  overtime_weekday_rate DECIMAL(4,2) DEFAULT 1.25,
  overtime_saturday_rate DECIMAL(4,2) DEFAULT 1.50,
  overtime_sunday_rate DECIMAL(4,2) DEFAULT 2.00,
  annual_leave_days INT DEFAULT 30,
  minimum_salary DECIMAL(15,2) DEFAULT 200000,
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Insert default company settings
INSERT INTO company_settings (
  company_name,
  company_address,
  company_nif,
  company_stat,
  cnaps_number,
  ostie_number,
  contact_email,
  contact_phone,
  timezone,
  currency,
  fiscal_year_start,
  work_hours_per_day,
  work_days_per_week,
  overtime_weekday_rate,
  overtime_saturday_rate,
  overtime_sunday_rate,
  annual_leave_days,
  minimum_salary
) VALUES (
  'PeopleDesk Madagascar',
  '123 Rue de l''Indépendance, Antananarivo',
  '5001234567',
  '12345678901234',
  '123456',
  '654321',
  'contact@peopledesk.mg',
  '+261 20 12 345 67',
  'Indian/Antananarivo',
  'MGA',
  '01-01',
  8.00,
  5,
  1.25,
  1.50,
  2.00,
  30,
  200000
);

-- Index for performance
CREATE INDEX idx_company_settings_updated_by ON company_settings(updated_by);