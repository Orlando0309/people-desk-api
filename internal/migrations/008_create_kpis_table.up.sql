-- KPIs table for performance management
CREATE TABLE kpis (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  target_value NUMERIC(15,2) NOT NULL,
  weight_percentage NUMERIC(5,2) NOT NULL CHECK (weight_percentage > 0 AND weight_percentage <= 100),
  scoring_scale VARCHAR(20) DEFAULT '1_to_5' CHECK (scoring_scale IN ('1_to_5', '1_to_10', 'custom')),
  department VARCHAR(100),
  position VARCHAR(100),
  is_active BOOLEAN DEFAULT true,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Performance Reviews table
CREATE TABLE performance_reviews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
  kpi_id UUID NOT NULL REFERENCES kpis(id) ON DELETE CASCADE,
  review_period_start DATE NOT NULL,
  review_period_end DATE NOT NULL,
  self_score NUMERIC(5,2),
  manager_score NUMERIC(5,2),
  final_score NUMERIC(5,2),
  self_assessment TEXT,
  manager_assessment TEXT,
  reviewer_id UUID NOT NULL REFERENCES users(id),
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'approved')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(employee_id, kpi_id, review_period_start, review_period_end)
);

-- Indexes for performance
CREATE INDEX idx_kpis_department ON kpis(department);
CREATE INDEX idx_kpis_position ON kpis(position);
CREATE INDEX idx_kpis_created_by ON kpis(created_by);
CREATE INDEX idx_performance_reviews_employee_id ON performance_reviews(employee_id);
CREATE INDEX idx_performance_reviews_kpi_id ON performance_reviews(kpi_id);
CREATE INDEX idx_performance_reviews_period ON performance_reviews(review_period_start, review_period_end);
CREATE INDEX idx_performance_reviews_status ON performance_reviews(status);