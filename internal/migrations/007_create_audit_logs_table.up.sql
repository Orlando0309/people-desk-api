-- Audit Logs table (CRITICAL: Immutable)
CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_role VARCHAR(20) NOT NULL,
  ip_address INET NOT NULL,
  action_type VARCHAR(50) NOT NULL,
  module VARCHAR(50) NOT NULL,
  record_id UUID,
  before_value JSONB,
  after_value JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_user_role ON audit_logs(user_role);
CREATE INDEX idx_audit_logs_action_type ON audit_logs(action_type);
CREATE INDEX idx_audit_logs_module ON audit_logs(module);
CREATE INDEX idx_audit_logs_record_id ON audit_logs(record_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_user_action ON audit_logs(user_id, action_type);
CREATE INDEX idx_audit_logs_module_created ON audit_logs(module, created_at);

-- Prevent deletion/modification of audit logs (immutable)
REVOKE DELETE, UPDATE ON audit_logs FROM PUBLIC;