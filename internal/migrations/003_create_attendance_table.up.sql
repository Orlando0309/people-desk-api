-- Attendance table
CREATE TABLE attendance (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  clock_in TIMESTAMPTZ,
  clock_out TIMESTAMPTZ,
  ip_address INET,
  device_fingerprint TEXT,
  status VARCHAR(20) DEFAULT 'present' CHECK (status IN ('present', 'absent', 'late', 'overtime', 'half_day')),
  total_hours NUMERIC(5,2),
  overtime_hours NUMERIC(5,2) DEFAULT 0,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  UNIQUE(employee_id, date)
);

-- Indexes for performance
CREATE INDEX idx_attendance_employee_id ON attendance(employee_id);
CREATE INDEX idx_attendance_date ON attendance(date);
CREATE INDEX idx_attendance_status ON attendance(status);
CREATE INDEX idx_attendance_deleted_at ON attendance(deleted_at);
CREATE INDEX idx_attendance_employee_date ON attendance(employee_id, date);