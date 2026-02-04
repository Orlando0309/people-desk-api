-- Employees table
CREATE TABLE employees (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  date_of_birth DATE,
  gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other')),
  nationality VARCHAR(100) DEFAULT 'Malagasy',
  national_id VARCHAR(50),
  position VARCHAR(100),
  department VARCHAR(100),
  hire_date DATE NOT NULL,
  contract_type VARCHAR(50) CHECK (contract_type IN ('permanent', 'fixed_term', 'intern', 'contractor')) DEFAULT 'permanent',
  gross_salary NUMERIC(15,2) NOT NULL CHECK (gross_salary >= 200000),
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'on_leave', 'terminated')),
  address TEXT,
  phone VARCHAR(50),
  emergency_contact_name VARCHAR(100),
  emergency_contact_phone VARCHAR(50),
  manager_id UUID REFERENCES employees(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Indexes for performance
CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_employees_first_name ON employees(first_name);
CREATE INDEX idx_employees_last_name ON employees(last_name);
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_employees_position ON employees(position);
CREATE INDEX idx_employees_status ON employees(status);
CREATE INDEX idx_employees_manager_id ON employees(manager_id);
CREATE INDEX idx_employees_deleted_at ON employees(deleted_at);
CREATE INDEX idx_employees_full_name ON employees(first_name, last_name);

-- Add foreign key constraint to users table
ALTER TABLE users ADD CONSTRAINT fk_users_employee_id FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE SET NULL;