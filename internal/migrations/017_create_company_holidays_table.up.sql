-- Create company_holidays table
CREATE TABLE IF NOT EXISTS company_holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Create index on date
CREATE INDEX IF NOT EXISTS idx_company_holidays_date ON company_holidays(date);

-- Create index on is_recurring
CREATE INDEX IF NOT EXISTS idx_company_holidays_is_recurring ON company_holidays(is_recurring);
