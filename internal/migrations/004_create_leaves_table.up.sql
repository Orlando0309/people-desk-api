-- Leaves table
CREATE TABLE leaves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
  leave_type VARCHAR(50) CHECK (leave_type IN ('annual', 'sick', 'maternity', 'exceptional', 'paternity', 'unpaid')) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  days_requested NUMERIC(5,2) NOT NULL,
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')) NOT NULL,
  approver_id UUID REFERENCES users(id) ON DELETE SET NULL,
  reason TEXT,
  rejection_reason TEXT,
  approved_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Indexes for performance
CREATE INDEX idx_leaves_employee_id ON leaves(employee_id);
CREATE INDEX idx_leaves_leave_type ON leaves(leave_type);
CREATE INDEX idx_leaves_status ON leaves(status);
CREATE INDEX idx_leaves_start_date ON leaves(start_date);
CREATE INDEX idx_leaves_end_date ON leaves(end_date);
CREATE INDEX idx_leaves_approver_id ON leaves(approver_id);
CREATE INDEX idx_leaves_deleted_at ON leaves(deleted_at);
CREATE INDEX idx_leaves_employee_status ON leaves(employee_id, status);