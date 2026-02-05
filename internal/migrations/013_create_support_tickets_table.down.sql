-- Drop indexes
DROP INDEX IF EXISTS idx_support_ticket_replies_created_at;
DROP INDEX IF EXISTS idx_support_ticket_replies_ticket_id;
DROP INDEX IF EXISTS idx_support_tickets_created_at;
DROP INDEX IF EXISTS idx_support_tickets_assigned_to;
DROP INDEX IF EXISTS idx_support_tickets_status;
DROP INDEX IF EXISTS idx_support_tickets_user_id;

-- Drop tables
DROP TABLE IF EXISTS support_ticket_replies;
DROP TABLE IF EXISTS support_tickets;