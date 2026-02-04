# Database Migration Tool

This tool manages database migrations for the People Desk API project.

## Prerequisites

- PostgreSQL database running
- Environment variables configured in `.env` file:
  - `DB_HOST` - Database host
  - `DB_PORT` - Database port
  - `DB_USER` - Database user
  - `DB_PASSWORD` - Database password
  - `DB_DATABASE` - Database name

## Usage

### Run all pending migrations
```bash
go run cmd/migrate/main.go up
```

### Run migrations up to a specific version
```bash
go run cmd/migrate/main.go up 005
```

### Rollback the last migration
```bash
go run cmd/migrate/main.go down
```

### Rollback to a specific version
```bash
go run cmd/migrate/main.go down 003
```

### Show migration status
```bash
go run cmd/migrate/main.go status
```

### Create a new migration
```bash
go run cmd/migrate/main.go create add_new_column
```

This creates two files:
- `internal/migrations/011_add_new_column.up.sql` - SQL to apply the migration
- `internal/migrations/011_add_new_column.down.sql` - SQL to rollback the migration

### Reset the database (rollback all and reapply)
```bash
go run cmd/migrate/main.go reset
```

## Migration File Naming Convention

Migration files must follow this pattern:
- Up migration: `{version}_{name}.up.sql`
- Down migration: `{version}_{name}.down.sql`

Example:
- `001_create_users_table.up.sql`
- `001_create_users_table.down.sql`

## How It Works

1. The tool reads all migration files from `internal/migrations/`
2. It creates a `schema_migrations` table to track applied migrations
3. Migrations are executed in version order (001, 002, 003, etc.)
4. Each migration is wrapped in a transaction for safety
5. The tool records which migrations have been applied

## Example Workflow

1. Create a new migration:
   ```bash
   go run cmd/migrate/main.go create add_email_index
   ```

2. Edit the generated files with your SQL:
   - `internal/migrations/011_add_email_index.up.sql`
   - `internal/migrations/011_add_email_index.down.sql`

3. Run the migration:
   ```bash
   go run cmd/migrate/main.go up
   ```

4. Check status:
   ```bash
   go run cmd/migrate/main.go status
   ```

5. If needed, rollback:
   ```bash
   go run cmd/migrate/main.go down
   ```
