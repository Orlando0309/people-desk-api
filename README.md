# PeopleDesk API

A comprehensive HR Management System API designed specifically for Madagascar businesses, ensuring compliance with Madagascar labor law, CNAPS (pension), OSTIE/OSIE (health insurance), and IRSA (income tax) regulations.

## 🌟 Features

### ✅ Implemented Modules

1. **Authentication & Authorization**
   - JWT-based authentication with refresh tokens
   - Role-based access control (Admin, HR, Accountant, Employee)
   - Password hashing with bcrypt
   - Account locking after failed login attempts

2. **Employee Management**
   - Full CRUD operations for employee records
   - Search and filtering capabilities
   - Organizational hierarchy (manager-subordinate relationships)
   - Department and position management
   - Minimum salary validation (200,000 MGA)

3. **Attendance Management**
   - Clock-in/clock-out tracking with IP address logging
   - Automatic late detection (after 9:00 AM)
   - Overtime calculation (> 8 hours per day)
   - Attendance statistics and reporting
   - Attendance correction workflow (HR/Admin only)

4. **Leave Management**
   - Multiple leave types (annual, sick, maternity, paternity, exceptional, unpaid)
   - Leave balance tracking (30 days/year per Madagascar law)
   - Approval workflow with rejection reasons
   - Overlap detection
   - Leave statistics per employee

5. **Audit Trail**
   - Immutable logging of all user actions
   - Comprehensive audit log filtering
   - Export to CSV for compliance audits
   - Record history tracking
   - Audit statistics dashboard

### 🚧 Pending Modules (To Be Implemented)

- **Payroll Draft Module** - HR salary preparation with CNAPS/OSTIE/IRSA calculations
- **Payroll Approval Module** - Accountant approval with GL entries
- **Declarations Module** - CNAPS, OSTIE, IRSA monthly forms
- **KPI & Performance Module** - Performance reviews and KPI tracking

## 🏗️ Technology Stack

- **Backend:** Go 1.21+ (Gin framework)
- **Database:** PostgreSQL 14+
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt/jwt)
- **Password Hashing:** bcrypt
- **Development:** Air (live reloading)

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/people-desk-api.git
cd people-desk-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_DATABASE=peopledesk
SERVER_PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

### 4. Set Up Database

Create the PostgreSQL database:

```bash
createdb peopledesk
```

Run migrations:

```bash
# Migrations will be auto-applied on first run
# Or manually run migration files in internal/migrations/
```

### 5. Run the Application

#### Development Mode (with live reloading):

```bash
air
```

#### Production Mode:

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Documentation

See [`API_DOCUMENTATION.md`](API_DOCUMENTATION.md) for complete API endpoint documentation.

### Quick Start Endpoints

- **Health Check:** `GET /health`
- **Login:** `POST /api/v1/auth/login`
- **Register User:** `POST /api/v1/auth/register` (Admin only)
- **List Employees:** `GET /api/v1/employees`
- **Clock In:** `POST /api/v1/attendance/clock-in`
- **Submit Leave:** `POST /api/v1/leaves`

## 🔐 Role-Based Access Control

### Roles:

1. **Admin** - Full system access including user management and audit logs
2. **HR** - Employee management, attendance tracking, leave approval
3. **Accountant** - Financial operations (payroll approval, declarations)
4. **Employee** - Self-service (own attendance, leave requests, profile)

### Permission Matrix:

| Feature | Admin | HR | Accountant | Employee |
|---------|-------|----|-----------| ---------|
| User Management | ✅ | ❌ | ❌ | ❌ |
| Employee CRUD | ✅ | ✅ | View only | Self only |
| Attendance Tracking | ✅ | ✅ | ✅ | ✅ Self |
| Attendance Correction | ✅ | ✅ | ❌ | ❌ |
| Leave Requests | ✅ | ✅ | ✅ | ✅ |
| Leave Approval | ✅ | ✅ | ❌ | ❌ |
| Audit Logs | ✅ | ❌ | ❌ | ❌ |

## 🗄️ Database Schema

The application uses the following main tables:

- `users` - User accounts with authentication
- `employees` - Employee records
- `attendance` - Attendance tracking
- `leaves` - Leave requests and approvals
- `payroll_drafts` - HR salary preparations
- `payroll_approved` - Accountant-approved payroll
- `audit_logs` - Immutable audit trail

See migration files in [`internal/migrations/`](internal/migrations/) for complete schema.

## 🇲🇬 Madagascar-Specific Features

### Salary Constraints
- Minimum salary: **200,000 MGA** (enforced at database level)

### Leave Balance
- Annual leave: **30 days per year** (2.5 days per month)
- Sick leave: Unlimited with medical certificate
- Maternity/Paternity leave: As per Madagascar labor law

### Attendance
- Late detection: After **9:00 AM**
- Overtime calculation: > **8 hours per day**
- Overtime rates (to be implemented in payroll):
  - Weekdays: 25% premium
  - Saturdays: 50% premium
  - Sundays/holidays: 100% premium

### Social Contributions (To Be Implemented in Payroll)
- **CNAPS** (Pension): 13% employer + 1% employee (capped at 1,600,000 MGA)
- **OSTIE** (Health): 5% employer + 1% employee (capped at 1,600,000 MGA)
- **IRSA** (Income Tax): Progressive scale (0-20%)

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 📁 Project Structure

```
people-desk-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── attendance/              # Attendance module
│   ├── audit/                   # Audit trail module
│   ├── auth/                    # Authentication module
│   ├── config/                  # Configuration management
│   ├── db/                      # Database connection
│   ├── employee/                # Employee management module
│   ├── leave/                   # Leave management module
│   ├── middleware/              # HTTP middlewares (auth, RBAC)
│   ├── migrations/              # Database migrations
│   ├── server/                  # HTTP server and routing
│   └── support/                 # Support/complaints module
├── .air.toml                    # Air configuration for live reloading
├── .env.example                 # Example environment variables
├── .gitignore                   # Git ignore rules
├── API_DOCUMENTATION.md         # Complete API documentation
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── prd.md                       # Product Requirements Document
└── README.md                    # This file
```

## 🔧 Development

### Live Reloading

The project uses Air for live reloading during development:

```bash
air
```

### Code Style

Follow Go best practices and conventions:
- Use `gofmt` for formatting
- Use `golint` for linting
- Write meaningful commit messages

### Adding New Modules

1. Create a new directory under `internal/`
2. Implement the following files:
   - `{module}_model.go` - Data models and DTOs
   - `{module}_repo.go` - Database operations
   - `{module}_handler.go` - HTTP handlers
   - `{module}_routes.go` - Route registration
3. Register routes in `internal/server/router.go`
4. Create migration files if needed

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is proprietary software. All rights reserved.

## 👥 Authors

- Development Team

## 🙏 Acknowledgments

- Madagascar labor law compliance requirements
- OHADA accounting standards
- CNAPS, OSTIE, and IRSA regulatory frameworks

## 📞 Support

For support, email support@peopledesk.mg or create an issue in the repository.

---

**Note:** This is an active development project. The Payroll, Declarations, and KPI modules are currently under development.
