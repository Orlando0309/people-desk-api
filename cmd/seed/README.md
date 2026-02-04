# Seed Command

This command is used to seed the database with default data for the People Desk API.

## Usage

```bash
go run cmd/seed/main.go <command>
```

## Available Commands

| Command | Description |
|---------|-------------|
| `all` | Seed all default data (users, employees, IRSA tax brackets, and payroll configurations) |
| `users` | Seed default users only |
| `employees` | Seed default employees only |
| `irsa` | Seed default IRSA tax brackets only |
| `payroll-config` | Seed default payroll configurations only |

## Examples

```bash
# Seed all default data
go run cmd/seed/main.go all

# Seed only default users
go run cmd/seed/main.go users

# Seed only default employees
go run cmd/seed/main.go employees

# Seed only IRSA tax brackets
go run cmd/seed/main.go irsa

# Seed only payroll configurations
go run cmd/seed/main.go payroll-config
```

## Default Data

### Default Users

The following default users are created:

| Email | Password | Role |
|-------|----------|------|
| `admin@peopledesk.com` | `admin123` | admin |
| `hr@peopledesk.com` | `hr123` | hr |
| `accountant@peopledesk.com` | `accountant123` | accountant |
| `employee@peopledesk.com` | `employee123` | employee |

**⚠️ Important:** These are default credentials for development/testing purposes. Change them in production!

### Default Employees

The following default employees are created:

1. **Admin User** - System Administrator (IT)
2. **HR Manager** - HR Manager (Human Resources)
3. **Accountant Specialist** - Senior Accountant (Finance)
4. **John Doe** - Software Developer (IT)
5. **Jane Smith** - Marketing Specialist (Marketing)

### Default IRSA Tax Brackets

The following IRSA tax brackets are seeded (example values for Madagascar):

| Bracket | Min Income (Ar) | Max Income (Ar) | Tax Rate | Min Tax (Ar) |
|---------|-----------------|-----------------|----------|--------------|
| Tranche 1 | 0 | 350,000 | 0% | 0 |
| Tranche 2 | 350,001 | 400,000 | 5% | 0 |
| Tranche 3 | 400,001 | 500,000 | 10% | 2,500 |
| Tranche 4 | 500,001 | 600,000 | 15% | 12,500 |
| Tranche 5 | 600,001 | ∞ | 20% | 27,500 |

**Note:** Adjust these values according to the actual IRSA regulations in Madagascar.

### Default Payroll Configurations

The following payroll configurations are seeded:

#### General Settings
- `minimum_wage`: 200,000 Ar
- `standard_working_hours`: 173.33 hours/month
- `payroll_processing_day`: 25
- `currency_code`: MGA
- `enable_auto_payroll`: false
- `rounding_method`: nearest

#### Overtime Settings
- `overtime_multiplier_1`: 1.30 (30% for first 8 hours)
- `overtime_multiplier_2`: 1.50 (50% for hours beyond 8)

#### Social Security Settings
- `cnaps_rate_employee`: 1%
- `cnaps_rate_employer`: 13%
- `ostie_rate_employee`: 1%
- `ostie_rate_employer`: 5%

#### Allowances
- `medical_allowance`: 50,000 Ar
- `transport_allowance`: 30,000 Ar
- `housing_allowance`: 100,000 Ar
- `family_allowance_rate`: 5%

#### Leave Settings
- `annual_leave_days`: 30 days
- `sick_leave_days`: 15 days

#### Tax Settings
- `tax_year_start`: 01-01
- `irsa_enabled`: true

## Prerequisites

Before running the seed command, ensure that:

1. The database migrations have been run:
   ```bash
   go run cmd/migrate/main.go up
   ```

2. The `.env` file is properly configured with database credentials.

## Notes

- The seed command checks if data already exists before creating it. Running the command multiple times will not create duplicates.
- All passwords are hashed using bcrypt.
- The seed command is intended for development and testing purposes. For production, consider using a proper data seeding strategy.