-- Add profile fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS phone VARCHAR(50),
ADD COLUMN IF NOT EXISTS address TEXT,
ADD COLUMN IF NOT EXISTS emergency_contact_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS emergency_contact_phone VARCHAR(50);

-- Add indexes for the new fields
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
