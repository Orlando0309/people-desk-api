# PeopleDesk API Documentation

## Base URL
```
http://localhost:{SERVER_PORT}/api/v1
```

## Authentication
All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer {access_token}
```

---

## 1. Authentication Endpoints

### POST /auth/login
Login to the system
- **Access:** Public
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
- **Response:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "admin",
    "is_active": true
  }
}
```

### POST /auth/refresh
Refresh access token
- **Access:** Public
- **Request Body:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

### POST /auth/register
Register a new user
- **Access:** Admin only
- **Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "role": "hr",
  "employee_id": "uuid (optional)"
}
```

### GET /auth/profile
Get current user profile
- **Access:** Authenticated users

### PUT /auth/change-password
Change password
- **Access:** Authenticated users
- **Request Body:**
```json
{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

### GET /auth/users
List all users
- **Access:** Admin only

### PUT /auth/users/:id
Update user
- **Access:** Admin only
- **Request Body:**
```json
{
  "email": "user@example.com",
  "role": "hr",
  "is_active": true
}
```

### DELETE /auth/users/:id
Delete user
- **Access:** Admin only

---

## 2. Dashboard Endpoints

### GET /dashboard/stats
Get dashboard statistics
- **Access:** Authenticated users
- **Response:**
```json
{
  "total_employees": 100,
  "on_leave_today": 5,
  "present_today": 90,
  "absent_today": 5,
  "overtime_hours_today": 12.5
}
```

### GET /dashboard/attendance/summary
Get attendance summary
- **Access:** Authenticated users
- **Response:**
```json
{
  "total_days": 20,
  "present_days": 18,
  "absent_days": 1,
  "late_days": 1,
  "average_hours": 8.5,
  "total_overtime": 25.5
}
```

### GET /dashboard/leaves/balances
Get leave balances for all employees
- **Access:** Authenticated users
- **Response:**
```json
{
  "leaves": [
    {
      "employee_id": "uuid",
      "employee_name": "John Doe",
      "annual_leave_balance": 15,
      "sick_leave_balance": 10,
      "maternity_leave_balance": 0,
      "paternity_leave_balance": 0
    }
  ]
}
```

---

## 3. Employee Management Endpoints

### GET /employees
List employees with filtering
- **Access:** All authenticated users
- **Query Parameters:**
  - `search` - Search by name or national ID
  - `department` - Filter by department
  - `position` - Filter by position
  - `status` - Filter by status (active, on_leave, terminated)
  - `limit` - Results per page (default: 50)
  - `offset` - Pagination offset

### POST /employees
Create a new employee
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "company_id": "uuid",
  "first_name": "John",
  "last_name": "Doe",
  "date_of_birth": "1990-01-01",
  "gender": "male",
  "nationality": "Malagasy",
  "national_id": "123456789",
  "position": "Developer",
  "department": "IT",
  "hire_date": "2024-01-01",
  "contract_type": "permanent",
  "gross_salary": 500000,
  "address": "123 Main St",
  "phone": "+261 34 12 345 67",
  "emergency_contact_name": "Jane Doe",
  "emergency_contact_phone": "+261 34 12 345 68",
  "manager_id": "uuid (optional)"
}
```

### GET /employees/:id
Get employee by ID
- **Access:** All authenticated users (employees can only view their own profile)

### PUT /employees/:id
Update employee
- **Access:** HR, Admin

### DELETE /employees/:id
Delete employee (soft delete)
- **Access:** Admin only

### GET /employees/departments
Get all unique departments
- **Access:** All authenticated users

### GET /employees/positions
Get all unique positions
- **Access:** All authenticated users

### GET /employees/:id/subordinates
Get employee's subordinates
- **Access:** All authenticated users

---

## 4. Attendance Management Endpoints

### POST /attendance/clock-in
Clock in for the day
- **Access:** All authenticated users
- **Request Body:**
```json
{
  "employee_id": "uuid"
}
```

### POST /attendance/clock-out
Clock out for the day
- **Access:** All authenticated users
- **Request Body:**
```json
{
  "employee_id": "uuid"
}
```

### GET /attendance/today/:employee_id
Get today's attendance for an employee
- **Access:** All authenticated users (employees can only view their own)

### GET /attendance
List attendance records
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `employee_id` - Filter by employee
  - `start_date` - Filter by start date
  - `end_date` - Filter by end date
  - `status` - Filter by status (present, absent, late, overtime, half_day)
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /attendance/:id
Get attendance record by ID
- **Access:** All authenticated users

### PUT /attendance/:id
Update attendance record (correction)
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "clock_in": "2024-01-01T08:00:00Z",
  "clock_out": "2024-01-01T17:00:00Z",
  "status": "present",
  "notes": "Correction reason"
}
```

### GET /attendance/stats/:employee_id
Get attendance statistics
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `start_date` - Start date for stats
  - `end_date` - End date for stats

---

## 5. Leave Management Endpoints

### GET /leaves
List leave requests
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `employee_id` - Filter by employee
  - `leave_type` - Filter by type (annual, sick, maternity, exceptional, paternity, unpaid)
  - `status` - Filter by status (pending, approved, rejected, cancelled)
  - `start_date` - Filter by start date
  - `end_date` - Filter by end date
  - `limit` - Results per page
  - `offset` - Pagination offset

### POST /leaves
Create a leave request
- **Access:** All authenticated users
- **Request Body:**
```json
{
  "employee_id": "uuid",
  "leave_type": "annual",
  "start_date": "2024-01-15",
  "end_date": "2024-01-20",
  "days_requested": 5,
  "reason": "Family vacation"
}
```

### GET /leaves/pending
Get all pending leave requests
- **Access:** HR, Admin

### GET /leaves/balance/:employee_id
Get leave balance for an employee
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `year` - Year for balance calculation (default: current year)

### GET /leaves/:id
Get leave request by ID
- **Access:** All authenticated users

### PUT /leaves/:id
Update leave request (only if pending)
- **Access:** All authenticated users (only their own)

### DELETE /leaves/:id
Cancel leave request
- **Access:** All authenticated users (only their own)

### PUT /leaves/:id/approve
Approve a leave request
- **Access:** HR, Admin

### PUT /leaves/:id/reject
Reject a leave request
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "approver_id": "uuid",
  "rejection_reason": "Insufficient staffing during requested period"
}
```

---

## 6. Audit Trail Endpoints

### GET /audit/logs
List audit logs with filtering
- **Access:** Admin only
- **Query Parameters:**
  - `user_id` - Filter by user
  - `user_role` - Filter by role (admin, hr, accountant, employee)
  - `action_type` - Filter by action type
  - `module` - Filter by module
  - `start_date` - Filter by start date
  - `end_date` - Filter by end date
  - `limit` - Results per page (default: 100)
  - `offset` - Pagination offset

### GET /audit/logs/export
Export audit logs to CSV
- **Access:** Admin only
- **Query Parameters:** Same as list audit logs

### GET /audit/logs/:id
Get audit log by ID
- **Access:** Admin only

### GET /audit/logs/user/:user_id
Get all audit logs for a specific user
- **Access:** Admin only
- **Query Parameters:**
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /audit/logs/module/:module
Get all audit logs for a specific module
- **Access:** Admin only
- **Query Parameters:**
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /audit/stats
Get audit statistics
- **Access:** Admin only
- **Query Parameters:**
  - `start_date` - Start date for stats
  - `end_date` - End date for stats

### GET /audit/history/:record_id
Get complete history for a specific record
- **Access:** Admin only

---

## 7. Payroll Module Endpoints

### POST /payroll/drafts
Create a new payroll draft
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "period_start": "2026-01-01",
  "period_end": "2026-01-31",
  "employee_id": "uuid",
  "gross_salary": 500000
}
```
- **Response:** Automatically calculates CNAPS (13%+1%), OSTIE (5%+1%), and IRSA based on Madagascar regulations

### GET /payroll/drafts
List payroll drafts
- **Access:** All authenticated users
- **Query Parameters:**
  - `period_start` - Filter by period start
  - `period_end` - Filter by period end
  - `employee_id` - Filter by employee
  - `status` - Filter by status (draft, approved)
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /payroll/drafts/:id
Get payroll draft by ID
- **Access:** All authenticated users

### PUT /payroll/drafts/:id
Update payroll draft
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "gross_salary": 550000
}
```

### DELETE /payroll/drafts/:id
Delete payroll draft
- **Access:** HR, Admin

### PUT /payroll/drafts/:id/approve
Approve payroll draft
- **Access:** Accountant, Admin
- Creates official fiche de paie with GL entries (OHADA compliant)
- Records digital signature from accountant

### GET /payroll/approved
List approved payrolls
- **Access:** All authenticated users
- **Query Parameters:**
  - `period_start` - Filter by period start
  - `period_end` - Filter by period end
  - `employee_id` - Filter by employee
  - `fiche_paie_number` - Filter by fiche paie number
  - `accountant_id` - Filter by accountant
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /payroll/approved/:id
Get approved payroll by ID
- **Access:** All authenticated users

### GET /payroll/approved/fiche/:fiche_paie_number
Get approved payroll by fiche paie number
- **Access:** All authenticated users

### GET /payroll/approved/:id/fiche-paie
Generate fiche de paie (payslip)
- **Access:** All authenticated users
- **Response:** Official payslip with digital signature

### GET /payroll/reconciliation
Generate reconciliation report
- **Access:** Accountant, Admin
- **Query Parameters:**
  - `start_date` - Start date for report (default: current month)
  - `end_date` - End date for report (default: current month)
- **Response:** Compares HR draft totals, accountant approved totals, and GL recorded amounts

---

## 8. KPI & Performance Management Endpoints

### POST /kpi
Create a new KPI template
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "name": "Sales Target",
  "description": "Monthly sales revenue target",
  "target_value": 1000000,
  "weight_percentage": 25,
  "scoring_scale": "1_to_5",
  "department": "Sales",
  "position": "Sales Manager"
}
```

### GET /kpi
List KPIs
- **Access:** All authenticated users
- **Query Parameters:**
  - `department` - Filter by department
  - `position` - Filter by position
  - `is_active` - Filter by active status
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /kpi/:id
Get KPI by ID
- **Access:** All authenticated users

### PUT /kpi/:id
Update KPI
- **Access:** HR, Admin

### DELETE /kpi/:id
Delete KPI
- **Access:** HR, Admin

### POST /kpi/reviews
Create a performance review
- **Access:** HR, Admin
- **Request Body:**
```json
{
  "employee_id": "uuid",
  "kpi_id": "uuid",
  "review_period_start": "2026-01-01",
  "review_period_end": "2026-03-31",
  "self_score": 4.5,
  "manager_score": 4.0,
  "self_assessment": "I met most targets...",
  "manager_assessment": "Good performance overall..."
}
```

### GET /kpi/reviews
List performance reviews
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `employee_id` - Filter by employee
  - `kpi_id` - Filter by KPI
  - `reviewer_id` - Filter by reviewer
  - `status` - Filter by status (pending, completed, approved)
  - `review_period_start` - Filter by period start
  - `review_period_end` - Filter by period end
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /kpi/reviews/:id
Get performance review by ID
- **Access:** All authenticated users (employees can only view their own)

### PUT /kpi/reviews/:id
Update performance review
- **Access:** All authenticated users
  - Employees: Can update self-assessment and self-score only
  - HR/Admin: Can update all fields

### DELETE /kpi/reviews/:id
Delete performance review
- **Access:** HR, Admin

### POST /kpi/reviews/:id/calculate
Calculate final score for performance review
- **Access:** HR, Admin
- Calculates weighted average of self and manager scores

### GET /kpi/reports/:employee_id
Generate performance report for an employee
- **Access:** All authenticated users (employees can only view their own)
- **Query Parameters:**
  - `start_date` - Start date for report (default: current quarter)
  - `end_date` - End date for report (default: current quarter)
- **Response:** Complete performance report with KPI breakdown and overall score

---

## 9. Declarations Module Endpoints

### POST /declarations
Create a new monthly declaration
- **Access:** Accountant, Admin
- **Request Body:**
```json
{
  "declaration_type": "cnaps",
  "declaration_period_start": "2026-01-01",
  "declaration_period_end": "2026-01-31",
  "company_name": "Company Name",
  "company_address": "Address",
  "company_nif": "123456789"
}
```

### GET /declarations
List declarations
- **Access:** All authenticated users
- **Query Parameters:**
  - `declaration_type` - Filter by type (cnaps, ostie, irsa)
  - `declaration_period_start` - Filter by period start
  - `declaration_period_end` - Filter by period end
  - `status` - Filter by status (draft, submitted, paid, cancelled)
  - `accountant_id` - Filter by accountant
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /declarations/:id
Get declaration by ID
- **Access:** All authenticated users

### GET /declarations/number/:declaration_number
Get declaration by declaration number
- **Access:** All authenticated users

### PUT /declarations/:id
Update declaration
- **Access:** Accountant, Admin
- **Request Body:**
```json
{
  "status": "submitted"
}
```

### DELETE /declarations/:id
Delete declaration
- **Access:** Accountant, Admin

### GET /declarations/:id/form
Generate declaration form
- **Access:** All authenticated users
- **Response:** Complete declaration form ready for submission

### POST /declarations/:id/populate
Populate declaration data from approved payrolls
- **Access:** Accountant, Admin
- Fills employee breakdown from payroll data

### GET /declarations/cnaps/generate
Generate CNAPS declaration form for a month
- **Access:** Accountant, Admin
- **Query Parameters:**
  - `month` - Month in format YYYY-MM (e.g., 2026-01)
- Creates or retrieves existing CNAPS declaration

### GET /declarations/ostie/generate
Generate OSTIE declaration form for a month
- **Access:** Accountant, Admin
- **Query Parameters:**
  - `month` - Month in format YYYY-MM (e.g., 2026-01)
- Creates or retrieves existing OSTIE declaration

### GET /declarations/irsa/generate
Generate IRSA declaration form for a month
- **Access:** Accountant, Admin
- **Query Parameters:**
  - `month` - Month in format YYYY-MM (e.g., 2026-01)
- Creates or retrieves existing IRSA declaration

### POST /declarations/irsa-brackets
Create IRSA tax bracket
- **Access:** Admin only
- **Request Body:**
```json
{
  "min_income": 0,
  "max_income": 350000,
  "tax_rate": 0,
  "min_tax": 2000,
  "effective_date": "2026-01-01"
}
```

### GET /declarations/irsa-brackets
List IRSA tax brackets
- **Access:** All authenticated users
- **Query Parameters:**
  - `is_active` - Filter by active status
  - `effective_date` - Filter by effective date
  - `limit` - Results per page
  - `offset` - Pagination offset

### GET /declarations/irsa-brackets/:id
Get IRSA tax bracket by ID
- **Access:** All authenticated users

### PUT /declarations/irsa-brackets/:id
Update IRSA tax bracket
- **Access:** Admin only

### DELETE /declarations/irsa-brackets/:id
Delete IRSA tax bracket
- **Access:** Admin only

---

## 10. Support/Complaints Endpoints

### POST /support
Submit a support ticket or complaint
- **Access:** All authenticated users
- **Request Body:**
```json
{
  "message": "Issue description",
  "email": "user@example.com"
}
```

---

## Role-Based Access Control (RBAC)

### Roles:
1. **admin** - Full system access
2. **hr** - HR operations (employee management, attendance, leave approval)
3. **accountant** - Financial operations (payroll approval, declarations)
4. **employee** - Self-service (own attendance, leave requests, profile)

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
| Payroll Draft | ✅ | ✅ | ❌ | ❌ |
| Payroll Approval | ✅ | ❌ | ✅ | ❌ |

---

## Error Responses

All endpoints return standard error responses:

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

---

## Madagascar-Specific Features

### Salary Constraints:
- Minimum salary: 200,000 MGA (enforced at database level)

### Leave Balance:
- Annual leave: 30 days per year (2.5 days per month)
- Sick leave: Unlimited with medical certificate
- Maternity leave: As per Madagascar labor law
- Paternity leave: As per Madagascar labor law

### Attendance:
- Late detection: After 9:00 AM
- Overtime calculation: > 8 hours per day
- Overtime rates (to be implemented in payroll):
  - Weekdays: 25% premium
  - Saturdays: 50% premium
  - Sundays/holidays: 100% premium

---

## Environment Variables Required

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_DATABASE=peopledesk
SERVER_PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

---

## Madagascar Payroll Calculations

### CNAPS (Caisse Nationale de Prévoyance Sociale)
- **Employee contribution:** 1% of gross salary (capped at 1,600,000 MGA)
- **Employer contribution:** 13% of gross salary (capped at 1,600,000 MGA)
- **Ceiling:** 8 × minimum wage (200,000 MGA) = 1,600,000 MGA

### OSTIE/OSIE (Organisme de Sécurité Sociale)
- **Employee contribution:** 1% of gross salary (capped at 1,600,000 MGA)
- **Employer contribution:** 5% of gross salary (capped at 1,600,000 MGA)
- **Ceiling:** Same as CNAPS (1,600,000 MGA)

### IRSA (Impôt sur les Revenus Salariaux et Assimilés)
Progressive tax scale on taxable income (gross - CNAPS employee - OSTIE employee):
- **≤ 350,000 MGA:** 0% (minimum tax: 2,000 MGA)
- **350,001–400,000 MGA:** 5%
- **400,001–500,000 MGA:** 10%
- **500,001–600,000 MGA:** 15%
- **> 600,000 MGA:** 20%

### OHADA-Compliant GL Entries
When payroll is approved, the following GL entries are created:
- **Account 641 (Salaires et traitements):** Debit gross salary
- **Account 421 (Salaires à payer):** Credit net salary
- **Account 431 (CNAPS à payer):** Credit CNAPS employee + employer contributions
- **Account 438 (OSTIE à payer):** Credit OSTIE employee + employer contributions
- **Account 437 (IRSA à verser):** Credit IRSA withholding
- **Account 646 (Charges sociales):** Debit employer CNAPS + employer OSTIE contributions

---

## Next Steps (To Be Implemented)

1. **Document Management** - Upload and manage employee documents
2. **Notifications** - Email/SMS notifications for leave approvals, etc.
3. **PDF Generation** - Generate downloadable PDF payslips and declaration forms
4. **Bank Integration** - Direct bank transfer for payroll payments
