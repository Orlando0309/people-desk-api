# **PRODUCT REQUIREMENTS DOCUMENT**  
## **PeopleDesk — Modular HR, Payroll & Accounting Platform for Madagascar**  

**Document Version:** 3.0  
**Date:** February 4, 2026  
**Product Owner:** [Client Name], Business Entrepreneur (Madagascar)  
**Target Market:** SMEs in Madagascar (10–500 employees)  
**Technology Stack:** Golang (Gin) + PostgreSQL + React (Latest) + Tailwind CSS  

---

## **TABLE OF CONTENTS**
- [1. EXECUTIVE SUMMARY](#1-executive-summary)
- [2. MADAGASCAR REGULATORY FRAMEWORK](#2-madagascar-regulatory-framework)
- [3. USER PERSONAS & ROLE-BASED ACCESS CONTROL](#3-user-personas--role-based-access-control)
- [4. PRODUCT SCOPE](#4-product-scope)
- [5. FUNCTIONAL REQUIREMENTS](#5-functional-requirements)
- [6. NON-FUNCTIONAL REQUIREMENTS](#6-non-functional-requirements)
- [7. TECHNICAL ARCHITECTURE](#7-technical-architecture)
- [8. ACCOUNTING INTEGRITY SAFEGUARDS](#8-accounting-integrity-safeguards)
- [9. SUCCESS METRICS](#9-success-metrics)
- [10. ASSUMPTIONS, CONSTRAINTS & RISKS](#10-assumptions-constraints--risks)
- [11. DEVELOPMENT ROADMAP](#11-development-roadmap)
- [12. OPEN QUESTIONS](#12-open-questions)
- [13. APPROVALS](#13-approvals)

---

## **1. EXECUTIVE SUMMARY**

**Product Vision:**  
PeopleDesk is a modular, web-based HR management platform designed specifically for Madagascar businesses. It digitizes end-to-end HR processes while ensuring strict compliance with Madagascar labor law, CNAPS (pension), OSTIE/OSIE (health insurance), and IRSA (income tax) regulations. The platform enforces separation of duties between HR and Accounting roles to prevent false accounting entries—a critical legal requirement under Madagascar law.

**Core Differentiators:**  
✅ Madagascar-specific regulatory compliance (CNAPS, OSTIE/OSIE, IRSA)  
✅ Strict role separation: HR prepares data → Accountant approves financial entries  
✅ Immutable audit trail for all user actions (Admin oversight)  
✅ OHADA-compliant accounting integration  
✅ Modular architecture: activate/deactivate features per client needs  

**Business Goal:**  
Acquire **3 paying companies within 6 months** of MVP launch.

---

## **2. MADAGASCAR REGULATORY FRAMEWORK**

### **2.1 Critical Distinctions: CNAPS vs. OSTIE/OSIE vs. IRSA**

| **Component** | **Nature** | **Employer Rate** | **Employee Rate** | **Ceiling** | **Purpose** |
|---------------|------------|-------------------|-------------------|-------------|-------------|
| **CNAPS**<br>(Caisse Nationale de Prévoyance Sociale) | Social security (pension) | 13% of gross salary | 1% of gross salary | 8 × minimum wage<br>(1,600,000 MGA)<sup>†</sup> | Funds employee retirement benefits |
| **OSTIE/OSIE**<br>(Organisme de Sécurité Sociale) | Statutory health insurance | 5% of gross salary | 1% of gross salary | Same ceiling as CNAPS<br>(1,600,000 MGA) | Covers medical expenses for employees & families (FUNHECE) |
| **IRSA**<br>(Impôt sur les Revenus Salariaux et Assimilés) | Withholding income tax | N/A (collection agent) | Progressive scale:<br>• ≤ 350,000 MGA: 0% (min. 2,000 MGA)<br>• 350,001–400,000: 5%<br>• 400,001–500,000: 10%<br>• 500,001–600,000: 15%<br>• > 600,000: 20% | No ceiling | Government tax revenue (remitted monthly to tax office) |

> <sup>†</sup> *Based on client-provided minimum wage of 200,000 MGA/month (2026). CNAPS/OSTIE ceiling = 8 × 200,000 = 1,600,000 MGA.*

### **2.2 Accounting Treatment (OHADA Compliance)**

| **Item** | **Correct Accounting Entry** | **Common Error to Prevent** |
|----------|------------------------------|----------------------------|
| **CNAPS Employer (13%)** | Liability account (431 – CNAPS à payer) | ❌ Recording as expense (641) |
| **CNAPS Employee (1%)** | Deduction from gross salary → Liability (431) | ❌ Treating as company expense |
| **OSTIE Employer (5%)** | Liability account (438 – OSTIE à payer) | ❌ Recording as direct expense |
| **OSTIE Employee (1%)** | Deduction from gross salary → Liability (438) | ❌ Misclassifying as benefit expense |
| **IRSA Withholding** | Liability account (437 – IRSA à verser) | ❌ Recording as company tax expense (631) |

> ⚠️ **Legal Warning:** Under Madagascar law, false accounting entries for social contributions constitute tax fraud. The employer acts as a *collection agent* for CNAPS/OSTIE/IRSA—not the beneficiary. All liabilities must be recorded separately and remitted monthly.

---

## **3. USER PERSONAS & ROLE-BASED ACCESS CONTROL**

### **3.1 User Roles**

| **Role** | **Description** | **Primary Responsibilities** |
|----------|-----------------|------------------------------|
| **Admin** | System administrator (typically company director/owner) | Full system configuration, module activation, user management, **audit trail review**, security settings |
| **HR Manager** | Human Resources professional | Employee lifecycle management, attendance tracking, leave approval, KPI management, **draft salary preparation** |
| **Accountant** | Certified accountant (comptable agréé) | **Sole authority** for financial entries, *fiche de paie* approval, CNAPS/OSTIE/IRSA declarations, general ledger management |
| **Employee** | Company staff member | View personal data, submit leave requests, clock in/out, view payslips, submit complaints |

### **3.2 Comprehensive RBAC Matrix**

| **Feature / Action** | **Admin** | **HR Manager** | **Accountant** | **Employee** |
|----------------------|-----------|----------------|----------------|--------------|
| **SYSTEM CONFIGURATION** | | | | |
| Activate/deactivate modules | ✅ Full control | ❌ | ❌ | ❌ |
| User role assignment | ✅ Full control | ❌ | ❌ | ❌ |
| Company settings (name, address, logo) | ✅ Full control | ❌ | ❌ | ❌ |
| Minimum wage / tax bracket configuration | ✅ Full control | ❌ | ✅ View only | ❌ |
| **AUDIT TRAIL** | | | | |
| View all user actions (who did what, when) | ✅ **Full visibility** | ❌ | ❌ | ❌ |
| Export audit logs (CSV/PDF) | ✅ | ❌ | ❌ | ❌ |
| Filter logs by user/date/action type | ✅ | ❌ | ❌ | ❌ |
| **EMPLOYEE MANAGEMENT** | | | | |
| Create employee profile | ✅ | ✅ | ❌ | ❌ |
| Edit personal data (name, address, etc.) | ✅ | ✅ | ❌ | ❌ |
| View all employee profiles | ✅ | ✅ | ✅ View only | ✅ Self only |
| Upload employee documents (contracts, diplomas) | ✅ | ✅ | ✅ View only | ✅ Self documents |
| Manage organizational chart | ✅ | ✅ | ✅ View only | ✅ View only |
| **ATTENDANCE MANAGEMENT** | | | | |
| Clock in/out (self) | ✅ | ✅ | ✅ | ✅ |
| View team attendance dashboard | ✅ | ✅ Full access | ✅ View only | ❌ Self only |
| Approve attendance corrections | ✅ | ✅ | ❌ | ❌ |
| Generate attendance reports | ✅ | ✅ | ✅ | ❌ |
| **LEAVE MANAGEMENT** | | | | |
| Submit leave request | ✅ | ✅ | ✅ | ✅ |
| View leave balance (self) | ✅ | ✅ | ✅ | ✅ |
| Approve/reject leave requests | ✅ | ✅ Full approval | ❌ | ❌ |
| Configure leave policies (days/year) | ✅ | ✅ | ❌ | ❌ |
| **KPI & PERFORMANCE** | | | | |
| Create KPI templates | ✅ | ✅ Full control | ❌ | ❌ |
| Conduct performance reviews | ✅ | ✅ | ❌ | ❌ View own only |
| Employee self-assessment | ✅ | ✅ | ✅ | ✅ Own only |
| **SALARY PREPARATION (HR DRAFT)** | | | | |
| Input gross salary per employee | ✅ | ✅ | ❌ | ❌ |
| Preview CNAPS calculation (13% + 1%) | ✅ | ✅ Preview only | ❌ | ❌ |
| Preview OSTIE calculation (5% + 1%) | ✅ | ✅ Preview only | ❌ | ❌ |
| Preview IRSA withholding | ✅ | ✅ Preview only | ❌ | ❌ |
| Generate *draft* fiche de paie | ✅ | ✅ Draft only<br>(watermarked "DRAFT") | ❌ | ❌ View draft only |
| **ACCOUNTING & PAYROLL (FINAL)** | | | | |
| **APPROVE salary calculations** | ✅ | ❌ **BLOCKED** | ✅ **SOLE AUTHORITY** | ❌ |
| Record CNAPS liabilities to GL | ✅ | ❌ **SYSTEM BLOCKED** | ✅ **MANDATORY** | ❌ |
| Record OSTIE liabilities to GL | ✅ | ❌ **SYSTEM BLOCKED** | ✅ **MANDATORY** | ❌ |
| Record IRSA withholding to GL | ✅ | ❌ **SYSTEM BLOCKED** | ✅ **MANDATORY** | ❌ |
| Generate official *fiche de paie* | ✅ | ❌ | ✅ Requires digital signature | ✅ View final only |
| Generate CNAPS monthly declaration | ✅ | ❌ | ✅ | ❌ |
| Generate OSTIE monthly declaration | ✅ | ❌ | ✅ | ❌ |
| Generate IRSA monthly declaration | ✅ | ❌ | ✅ | ❌ |
| Manage supplier invoices & payments | ✅ | ❌ | ✅ Full control | ❌ |
| Access general ledger (accounts 431, 437, 438) | ✅ | ❌ **BLOCKED** | ✅ Full access | ❌ |
| Financial reporting (P&L, balance sheet) | ✅ | ❌ | ✅ | ❌ |
| **COMPLAINTS & FEEDBACK** | | | | |
| Submit complaint/feedback | ✅ | ✅ | ✅ | ✅ |
| View team complaints | ✅ | ✅ HR oversight | ❌ | ❌ Self only |
| Resolve complaints | ✅ | ✅ | ❌ | ❌ |
| **REPORTING** | | | | |
| HR analytics dashboard | ✅ | ✅ | ✅ View only | ❌ |
| Financial/compliance reports | ✅ | ❌ | ✅ | ❌ |
| Export to Excel/PDF | ✅ All modules | ✅ HR modules only | ✅ Financial modules only | ✅ Self data only |

---

## **4. PRODUCT SCOPE**

### **4.1 In-Scope (MVP – Phase 1: Web Desktop Application)**

#### **Core Modules**
- **Employee Management**
  - Full employee lifecycle (onboarding → offboarding)
  - Document management (contracts, diplomas, ID cards)
  - Organizational chart with drag-and-drop
  - Search/filter by department, position, status

- **Attendance Management**
  - Web-based clock-in/clock-out with timestamp/IP logging
  - Attendance calendar with color-coded status (present/absent/late)
  - Overtime calculation (per Madagascar labor code)
  - Absence justification workflow

- **Leave Management**
  - Multi-type leave (annual, sick, maternity, exceptional)
  - Leave balance tracking with accrual rules
  - Approval workflow with escalation rules
  - Leave calendar visualization

- **KPI & Performance Management**
  - Customizable KPI templates per role/department
  - Quarterly/annual review cycles
  - 360° feedback collection
  - Performance history tracking

- **Salary Preparation (HR Draft)**
  - Gross salary input per employee
  - Automatic CNAPS calculation (13% employer + 1% employee, capped at 1,600,000 MGA)
  - Automatic OSTIE calculation (5% employer + 1% employee, capped at 1,600,000 MGA)
  - Automatic IRSA calculation (progressive scale)
  - Draft *fiche de paie* generation (watermarked "DRAFT")

- **Accounting & Payroll (Accountant Approval)**
  - Mandatory accountant approval workflow before finalization
  - OHADA-compliant GL entries:
    - Account 431: CNAPS à payer
    - Account 438: OSTIE à payer
    - Account 437: IRSA à verser
    - Account 641: Salaires et traitements (gross salary expense)
  - Official *fiche de paie* generation with accountant digital signature
  - Monthly declaration forms:
    - CNAPS Form (employer + employee contributions)
    - OSTIE Form (employer + employee contributions)
    - IRSA Form (withholding summary)

- **Complaints Management**
  - Anonymous or identified complaint submission
  - HR assignment and resolution tracking
  - Status updates visible to complainant

- **Audit Trail (Admin Only)**
  - Immutable log of all user actions:
    - User ID, role, timestamp, IP address
    - Action type (create/update/delete/approve)
    - Module affected (employee/attendance/leave/etc.)
    - Before/after values for data changes
  - Searchable/filterable audit interface
  - Export to CSV/PDF for compliance audits

#### **Technical Scope**
- Web application (responsive desktop-first design)
- Backend: Golang (Gin framework)+ GORM + PostgreSQL
- Authentication: JWT with refresh tokens
- Development: Air for live reloading
- Frontend: React 18+ with TypeScript
- Styling: Tailwind CSS + Headless UI
- Currency: MGA (Ariary) only – no multi-currency support

### **4.2 Out-of-Scope (Phase 1)**
- Mobile application (planned for v2)
- Payroll payment integration (bank transfers)
- Multi-language support (French/Malagasy only)
- Advanced analytics/AI recommendations
- API for third-party ERP/accounting systems
- Recruitment module
- Training & development module

---

## **5. FUNCTIONAL REQUIREMENTS**

### **5.1 Authentication & Authorization**
- **FR-AUTH-001:** JWT-based authentication with 15-minute access token + 7-day refresh token
- **FR-AUTH-002:** Role-based access control enforced at API gateway level (Gin middleware)
- **FR-AUTH-003:** Password policy: minimum 8 characters, 1 uppercase, 1 number, 1 special character
- **FR-AUTH-004:** Session timeout after 30 minutes of inactivity
- **FR-AUTH-005:** Failed login lockout after 5 attempts (15-minute cooldown)

### **5.2 Audit Trail (Critical Requirement)**
- **FR-AUDIT-001:** System SHALL log every user action with:
  - `user_id`, `user_role`, `timestamp`, `ip_address`, `action_type`, `module`, `record_id`, `before_value`, `after_value`
- **FR-AUDIT-002:** Audit logs SHALL be immutable (no DELETE/UPDATE operations allowed on audit table)
- **FR-AUDIT-003:** Admin SHALL access audit interface with filters:
  - Date range picker
  - User selector (single/multiple)
  - Action type (create/update/delete/approve/reject)
  - Module selector (employee/attendance/leave/salary/etc.)
- **FR-AUDIT-004:** Audit logs SHALL be retained for minimum 10 years (Madagascar legal requirement)
- **FR-AUDIT-005:** System SHALL generate daily reconciliation report:
  ```
  [HR Draft Totals] vs [Accountant Approved Totals] vs [GL Recorded Amounts]
  ```
  Discrepancies > 0.1% trigger email alert to Admin + Accountant

### **5.3 Employee Management**
- **FR-EMP-001:** Create employee profile with mandatory fields:
  - Personal: name, date of birth, gender, nationality, ID number
  - Contact: address, phone, email
  - Employment: position, department, hire date, contract type, gross salary
- **FR-EMP-002:** Upload documents with categorization (contract, diploma, ID, medical certificate)
- **FR-EMP-003:** Search employees by name, ID, department, position, status (active/inactive)
- **FR-EMP-004:** Generate organizational chart with reporting lines

### **5.4 Attendance Management**
- **FR-ATT-001:** Clock-in/clock-out with timestamp, IP address, and device fingerprint
- **FR-ATT-002:** Attendance calendar showing daily status (green = present, red = absent, yellow = late)
- **FR-ATT-003:** Overtime calculation:
  - Weekdays > 8h/day = 25% premium
  - Saturdays = 50% premium
  - Sundays/holidays = 100% premium (per Madagascar labor code)
- **FR-ATT-004:** Absence justification workflow (employee submits reason → HR approves/rejects)

### **5.5 Leave Management**
- **FR-LEAVE-001:** Submit leave request with:
  - Leave type (annual/sick/maternity/exceptional)
  - Start/end date + half-day option
  - Reason/justification
- **FR-LEAVE-002:** Automatic leave balance calculation:
  - Annual leave: 2.5 days per month worked (30 days/year max)
  - Sick leave: unlimited with medical certificate
- **FR-LEAVE-003:** Approval workflow:
  - Employee submits → HR approves/rejects with comment → Notification to employee
- **FR-LEAVE-004:** Leave calendar showing team availability

### **5.6 KPI & Performance Management**
- **FR-KPI-001:** Create KPI templates with:
  - Metric name, description, target value, weight (%)
  - Scoring scale (1–5 or custom)
- **FR-KPI-002:** Assign KPIs to employees/teams per review cycle
- **FR-KPI-003:** Employee self-assessment + manager evaluation
- **FR-KPI-004:** Performance score calculation (weighted average)
- **FR-KPI-005:** Generate performance reports (PDF/Excel)

### **5.7 Salary Preparation Workflow (Critical Safeguard)**
```
STEP 1: HR inputs gross salary per employee
  → System calculates:
     • CNAPS employee (1%) + employer (13%) [capped at 1,600,000 MGA]
     • OSTIE employee (1%) + employer (5%) [capped at 1,600,000 MGA]
     • IRSA withholding (progressive scale)
  → Generates DRAFT fiche de paie (watermarked "DRAFT - NOT OFFICIAL")

STEP 2: Accountant reviews draft
  → System displays reconciliation report:
     [HR Draft Total] vs [Expected Total per OHADA rules]
  → Accountant approves/rejects with comment

STEP 3: Upon approval
  → System records GL entries (accountant role ONLY):
     • Debit 641 (Salaires) | Credit 421 (Salaires à payer)
     • Debit 641 (Salaires) | Credit 431 (CNAPS à payer) [employee portion]
     • Debit 646 (Charges sociales) | Credit 431 (CNAPS à payer) [employer portion]
     • Debit 641 (Salaires) | Credit 438 (OSTIE à payer) [employee portion]
     • Debit 646 (Charges sociales) | Credit 438 (OSTIE à payer) [employer portion]
     • Debit 641 (Salaires) | Credit 437 (IRSA à verser)
  → Generates OFFICIAL fiche de paie with:
     • Accountant digital signature (JWT-verified)
     • Unique declaration number
     • Timestamp of approval
  → Locks record (immutable after approval)
```

- **FR-SALARY-001:** HR SHALL NOT access GL entry screens for accounts 431/437/438
- **FR-SALARY-002:** Backend validation SHALL reject GL entries attempted by non-accountant roles
- **FR-SALARY-003:** Official *fiche de paie* SHALL display mandatory fields:
  ```
  • Employee name, ID, position, department
  • Payment period (du [start] au [end])
  • Gross salary (Salaire Brut)
  • CNAPS employee deduction (1%)
  • OSTIE employee deduction (1%)
  • IRSA withholding (amount + bracket applied)
  • Net salary (Salaire Net à Payer)
  • Employer CNAPS (13%) + OSTIE (5%) for declaration purposes
  • Accountant name + digital signature timestamp
  • Unique declaration number (e.g., FDPAIE-2026-00123)
  ```

### **5.8 CNAPS/OSTIE/IRSA Declarations**
- **FR-DECL-001:** Accountant generates monthly CNAPS declaration form with:
  - Company identification (name, address, CNAPS number)
  - Employee list with gross salaries
  - CNAPS base (capped at 1,600,000 MGA)
  - Employee contribution (1%) + Employer contribution (13%)
  - Total employer contribution amount due
- **FR-DECL-002:** Accountant generates monthly OSTIE declaration form with:
  - Same structure as CNAPS but for health insurance (1% employee + 5% employer)
- **FR-DECL-003:** Accountant generates monthly IRSA declaration form with:
  - Employee taxable income after exemptions
  - IRSA calculated per progressive scale
  - Total IRSA withheld for remittance
- **FR-DECL-004:** System prevents duplicate declarations for same month/employee

### **5.9 Complaints Management**
- **FR-COMP-001:** Employee submits complaint with:
  - Anonymous/identified toggle
  - Category (harassment, discrimination, working conditions, etc.)
  - Description + optional attachment
- **FR-COMP-002:** HR assigned complaint with status tracking (New → In Progress → Resolved)
- **FR-COMP-003:** Complainant receives status updates via in-app notification

---

## **6. NON-FUNCTIONAL REQUIREMENTS**

### **6.1 Performance**
- **NFR-PERF-001:** Initial page load < 3 seconds on 4G connection
- **NFR-PERF-002:** API response time < 500ms for 95% of requests
- **NFR-PERF-003:** Support 500 concurrent users with < 2s response time
- **NFR-PERF-004:** Handle 10,000+ employee records with pagination/infinite scroll

### **6.2 Security**
- **NFR-SEC-001:** All data encrypted in transit (TLS 1.3) and at rest (AES-256)
- **NFR-SEC-002:** Passwords hashed with bcrypt (cost factor 12)
- **NFR-SEC-003:** SQL injection prevention via parameterized queries (GORM)
- **NFR-SEC-004:** XSS prevention via React DOM sanitization + Content Security Policy
- **NFR-SEC-005:** Role-based API middleware (Gin) validates permissions on every request
- **NFR-SEC-006:** Audit logs stored in separate PostgreSQL schema with restricted access

### **6.3 Usability**
- **NFR-USE-001:** Intuitive interface requiring < 30 minutes training for new users
- **NFR-USE-002:** Responsive design supporting 1280px+ desktop screens (mobile v2)
- **NFR-USE-003:** Consistent navigation with persistent sidebar menu
- **NFR-USE-004:** WCAG 2.1 AA compliance (color contrast, keyboard navigation)

### **6.4 Reliability**
- **NFR-REL-001:** 99.5% uptime SLA (excluding scheduled maintenance)
- **NFR-REL-002:** Automated daily backups with 7-day retention + weekly offsite backups
- **NFR-REL-003:** Error tracking via Sentry with real-time alerts
- **NFR-REL-004:** Database replication for disaster recovery

### **6.5 Accounting Compliance (Critical)**
- **NFR-ACC-001:** System SHALL prevent HR users from accessing GL accounts 431/437/438
- **NFR-ACC-002:** All financial entries SHALL comply with OHADA SYSCOHADA accounting plan
- **NFR-ACC-003:** Audit trail SHALL be tamper-proof (cryptographic hashing of log entries)
- **NFR-ACC-004:** System SHALL validate CNAPS/OSTIE calculations against ceiling (1,600,000 MGA)
- **NFR-ACC-005:** IRSA calculations SHALL use current progressive scale (configurable by Admin)

---

## **7. TECHNICAL ARCHITECTURE**

### **7.1 System Diagram**
```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT (React + Tailwind)                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   HR View   │  │ Accountant  │  │      Admin          │  │
│  │             │  │    View     │  │   (Audit Trail)     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└───────────────┬───────────────────────────┬─────────────────┘
                │        HTTPS (JWT)        │
                ▼                           ▼
┌─────────────────────────────────────────────────────────────┐
│               API GATEWAY (Gin Middleware)                  │
│  • Authentication (JWT validation)                          │
│  • Role-based access control                                │
│  • Rate limiting                                            │
│  • Request logging                                          │
└───────────────┬─────────────────────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────┐
│                  APPLICATION LAYERS                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │ Employee │  │Attendance│  │  Leave   │  │  Payroll   │  │
│  │  Service │  │ Service  │  │ Service  │  │  Service   │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────┬─────┘  │
│                                                  │         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐        │         │
│  │   KPI    │  │Complaints│  │  Audit   │◄───────┘         │
│  │ Service  │  │ Service  │  │ Service  │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└───────────────┬─────────────────────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────┐
│                    POSTGRESQL DATABASE                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │ employees│  │attendance│  │  leaves  │  │   payroll  │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │   kpis   │  │complaints│  │  users   │  │   audit    │  │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### **7.2 Core Database Schema (PostgreSQL)**

```sql
-- Users & Roles
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  role VARCHAR(20) CHECK (role IN ('admin', 'hr', 'accountant', 'employee')) NOT NULL,
  employee_id UUID REFERENCES employees(id),
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  last_login TIMESTAMPTZ
);

-- Employees
CREATE TABLE employees (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  company_id UUID NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  date_of_birth DATE,
  national_id VARCHAR(50),
  position VARCHAR(100),
  department VARCHAR(100),
  hire_date DATE NOT NULL,
  gross_salary NUMERIC(15,2) NOT NULL CHECK (gross_salary >= 200000),
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'on_leave', 'terminated')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Attendance
CREATE TABLE attendance (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID REFERENCES employees(id) NOT NULL,
  date DATE NOT NULL,
  clock_in TIMESTAMPTZ,
  clock_out TIMESTAMPTZ,
  ip_address INET,
  status VARCHAR(20) DEFAULT 'present' CHECK (status IN ('present', 'absent', 'late', 'overtime')),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Leaves
CREATE TABLE leaves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  employee_id UUID REFERENCES employees(id) NOT NULL,
  leave_type VARCHAR(50) CHECK (leave_type IN ('annual', 'sick', 'maternity', 'exceptional')) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')) NOT NULL,
  approver_id UUID REFERENCES users(id),
  reason TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Payroll (Draft by HR → Approved by Accountant)
CREATE TABLE payroll_drafts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  employee_id UUID REFERENCES employees(id) NOT NULL,
  gross_salary NUMERIC(15,2) NOT NULL,
  cnaps_employee NUMERIC(15,2) GENERATED ALWAYS AS (LEAST(gross_salary, 1600000) * 0.01) STORED,
  cnaps_employer NUMERIC(15,2) GENERATED ALWAYS AS (LEAST(gross_salary, 1600000) * 0.13) STORED,
  ostie_employee NUMERIC(15,2) GENERATED ALWAYS AS (LEAST(gross_salary, 1600000) * 0.01) STORED,
  ostie_employer NUMERIC(15,2) GENERATED ALWAYS AS (LEAST(gross_salary, 1600000) * 0.05) STORED,
  irsa NUMERIC(15,2) NOT NULL, -- Calculated via function
  net_salary NUMERIC(15,2) GENERATED ALWAYS AS (
    gross_salary - cnaps_employee - ostie_employee - irsa
  ) STORED,
  created_by UUID REFERENCES users(id) NOT NULL CHECK ((SELECT role FROM users WHERE id = created_by) = 'hr'),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE payroll_approved (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  draft_id UUID REFERENCES payroll_drafts(id) NOT NULL UNIQUE,
  fiche_paie_number VARCHAR(50) UNIQUE NOT NULL, -- e.g., FDPAIE-2026-00123
  accountant_id UUID REFERENCES users(id) NOT NULL CHECK ((SELECT role FROM users WHERE id = accountant_id) = 'accountant'),
  gl_entries JSONB NOT NULL, -- OHADA-compliant entries
  approved_at TIMESTAMPTZ DEFAULT NOW(),
  digital_signature TEXT NOT NULL -- JWT-based signature
);

-- CRITICAL: Audit Trail (Immutable)
CREATE TABLE audit_logs (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID REFERENCES users(id) NOT NULL,
  user_role VARCHAR(20) NOT NULL,
  ip_address INET NOT NULL,
  action_type VARCHAR(50) NOT NULL, -- 'create', 'update', 'delete', 'approve', etc.
  module VARCHAR(50) NOT NULL, -- 'employee', 'attendance', 'leave', 'payroll', etc.
  record_id UUID,
  before_value JSONB,
  after_value JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Prevent deletion/modification of audit logs
REVOKE DELETE, UPDATE ON audit_logs FROM PUBLIC;
CREATE POLICY audit_immutable ON audit_logs FOR ALL USING (true) WITH CHECK (false);
```

### **7.3 API Endpoints (Sample)**

```http
# Authentication
POST /api/auth/login          → { email, password } → { access_token, refresh_token }
POST /api/auth/refresh        → { refresh_token } → { access_token }

# Audit Trail (Admin only)
GET /api/audit/logs           → Query params: ?start_date=2026-01-01&end_date=2026-01-31&user_id=xxx&action_type=update
GET /api/audit/export         → CSV/PDF export of filtered logs

# Employees
GET /api/employees            → List (filtered by role permissions)
POST /api/employees           → Create (HR/Admin only)
PUT /api/employees/:id        → Update (HR/Admin only)
DELETE /api/employees/:id     → Soft delete (Admin only)

# Attendance
POST /api/attendance/clock-in → Employee self-service
POST /api/attendance/clock-out→ Employee self-service
GET /api/attendance/team      → HR view of team attendance

# Leave
POST /api/leaves              → Submit request (all roles)
PUT /api/leaves/:id/approve   → Approve (HR/Admin only)
PUT /api/leaves/:id/reject    → Reject (HR/Admin only)

# Payroll Draft (HR only)
POST /api/payroll/drafts      → Create draft salary calculations
GET /api/payroll/drafts       → List drafts awaiting approval

# Payroll Approval (Accountant only)
PUT /api/payroll/drafts/:id/approve → Approve draft → creates payroll_approved record + GL entries
GET /api/payroll/official     → List official fiches de paie

# Declarations (Accountant only)
GET /api/declarations/cnaps?month=2026-01 → Generate CNAPS declaration form
GET /api/declarations/ostie?month=2026-01 → Generate OSTIE declaration form
GET /api/declarations/irsa?month=2026-01  → Generate IRSA declaration form
```

---

## **8. ACCOUNTING INTEGRITY SAFEGUARDS**

### **8.1 Separation of Duties Enforcement**

| **Process Step** | **HR Role** | **Accountant Role** | **System Enforcement** |
|------------------|-------------|---------------------|------------------------|
| Input gross salary | ✅ Allowed | ❌ Blocked | API middleware validates role = 'hr' |
| Preview CNAPS/OSTIE/IRSA | ✅ Allowed (read-only) | ✅ Allowed | Calculations shown but no GL access |
| Approve salary draft | ❌ **BLOCKED** | ✅ **REQUIRED** | Backend validation: `if role != 'accountant' → 403 Forbidden` |
| Record CNAPS liability (GL 431) | ❌ **SYSTEM BLOCKED** | ✅ Allowed | PostgreSQL RLS policy: `CREATE POLICY cnaps_gl_policy ON gl_entries FOR INSERT TO accountant USING (account_code = '431')` |
| Record OSTIE liability (GL 438) | ❌ **SYSTEM BLOCKED** | ✅ Allowed | Same RLS enforcement |
| Record IRSA liability (GL 437) | ❌ **SYSTEM BLOCKED** | ✅ Allowed | Same RLS enforcement |
| Generate official *fiche de paie* | ❌ Blocked | ✅ Requires JWT signature | PDF generation requires accountant JWT token validation |

### **8.2 Reconciliation & Fraud Prevention**

- **Daily Reconciliation Report** (auto-generated at 02:00 AM):
  ```
  PAYROLL RECONCILIATION REPORT - 2026-02-04
  ===========================================
  HR Draft Totals (all employees):
    Gross Salary:       24,500,000 MGA
    CNAPS Employee:        245,000 MGA
    CNAPS Employer:      3,185,000 MGA
    OSTIE Employee:        245,000 MGA
    OSTIE Employer:      1,225,000 MGA
    IRSA Withheld:         850,000 MGA
    Net Payable:        22,860,000 MGA
  
  Accountant Approved Totals:
    Gross Salary:       24,500,000 MGA ✅
    CNAPS Employee:        245,000 MGA ✅
    CNAPS Employer:      3,185,000 MGA ✅
    OSTIE Employee:        245,000 MGA ✅
    OSTIE Employer:      1,225,000 MGA ✅
    IRSA Withheld:         850,000 MGA ✅
    Net Payable:        22,860,000 MGA ✅
  
  GL Recorded Amounts:
    Account 431 (CNAPS): 3,430,000 MGA ✅
    Account 438 (OSTIE): 1,470,000 MGA ✅
    Account 437 (IRSA):    850,000 MGA ✅
  
  STATUS: ✅ RECONCILED (0.0% variance)
  ```
  
- **Variance Alert Threshold:** > 0.1% triggers email to Admin + Accountant
- **Immutable Records:** Approved payroll records cannot be modified/deleted (soft delete only with audit trail)

---

## **9. SUCCESS METRICS**

### **9.1 Business KPIs**
| **Metric** | **Target (6 months)** | **Target (12 months)** |
|------------|------------------------|------------------------|
| Active paying companies | 3 | 10 |
| Monthly Recurring Revenue (MRR) | 1,500,000 MGA | 5,000,000 MGA |
| Employee users onboarded | 150 | 500 |
| Customer satisfaction (NPS) | ≥ 40 | ≥ 60 |

### **9.2 Product KPIs**
| **Metric** | **Target** |
|------------|------------|
| Daily Active Users (DAU) | ≥ 30% of total users |
| Feature adoption rate | ≥ 70% for core modules |
| Time to first *fiche de paie* generation | < 45 minutes |
| Support tickets related to accounting errors | 0 (critical metric) |

### **9.3 Technical KPIs**
| **Metric** | **Target** |
|------------|------------|
| Uptime | 99.5% |
| API error rate | < 0.5% |
| Page load time | < 3 seconds |
| Security vulnerabilities | 0 critical/high |

---

## **10. ASSUMPTIONS, CONSTRAINTS & RISKS**

### **10.1 Assumptions**
- Minimum wage remains 200,000 MGA for 2026 (client confirmation)
- Target clients have stable internet connectivity (4G/fiber)
- Accountants using the system hold valid Madagascar accounting certifications
- Companies operate under OHADA accounting standards

### **10.2 Constraints**
| **Constraint** | **Impact** |
|----------------|------------|
| Technology stack fixed (Golang + React) | Limits third-party library choices |
| Timeline: MVP in 4–5 months | Requires phased feature delivery |
| Budget: Competitive pricing for Madagascar SMEs | Must control hosting/development costs |
| Regulatory: Must comply with Madagascar labor law | Requires legal review before launch |

### **10.3 Risk Register**
| **Risk** | **Impact** | **Probability** | **Mitigation** |
|----------|------------|-----------------|----------------|
| HR bypasses accountant approval | **CRITICAL** – Legal liability | Medium | System-level RBAC enforcement + audit trail |
| CNAPS/OSTIE miscalculation due to outdated ceiling | High – Penalties from authorities | Medium | Admin-configurable ceiling + system alerts on outdated values |
| IRSA bracket changes mid-year | Medium – Employee disputes | Low | Admin-configurable tax tables with version history |
| Data breach exposing employee salaries | High – Reputational damage | Low | Encryption at rest/transit + strict access controls |
| Low user adoption due to complexity | Medium – Revenue shortfall | Medium | Extensive UX testing + onboarding tutorials |

---

## **11. DEVELOPMENT ROADMAP**

### **Phase 1: MVP Core (Months 1–4)**
| **Sprint** | **Deliverables** |
|------------|------------------|
| Sprint 1 | Auth system + User management + RBAC middleware |
| Sprint 2 | Employee management module + Document storage |
| Sprint 3 | Attendance tracking + Leave management |
| Sprint 4 | KPI module + Complaints system |
| Sprint 5 | Salary draft workflow (HR) + CNAPS/OSTIE/IRSA calculations |
| Sprint 6 | Accountant approval workflow + GL entries + *Fiche de paie* generation |
| Sprint 7 | Audit trail system + Admin dashboard |
| Sprint 8 | OHADA-compliant declarations + Testing + Security audit |

### **Phase 2: Polish & Scale (Months 5–6)**
- Performance optimization
- User acceptance testing with 3 pilot companies
- Regulatory compliance review (Madagascar labor law)
- Documentation & training materials

### **Phase 3: Mobile App (Months 7–9)**
- Employee mobile app (React Native)
  - Clock in/out
  - Leave requests
  - View *fiche de paie*
  - Submit complaints

### **Phase 4: Advanced Features (Months 10–12)**
- Recruitment module
- Training & development tracking
- Advanced analytics dashboard
- API for accounting software integration

---

## **12. OPEN QUESTIONS**

| **#** | **Question** | **Owner** | **Deadline** |
|-------|--------------|-----------|--------------|
| Q1 | Confirm exact CNAPS/OSTIE ceiling calculation methodology per latest Madagascar decree | Client | Feb 15, 2026 |
| Q2 | Provide sample *fiche de paie* format used by target clients for UI validation | Client | Feb 15, 2026 |
| Q3 | Confirm OHADA chart of accounts version (SYSCOHADA 2018 or 2024) | Accountant | Feb 20, 2026 |
| Q4 | Define pricing model (per employee/month vs. flat fee) | Client | Feb 25, 2026 |
| Q5 | Identify 3 pilot companies for UAT | Client | Mar 1, 2026 |

---

## **13. APPROVALS**

| **Role** | **Name** | **Signature** | **Date** | **Comments** |
|----------|----------|---------------|----------|--------------|
| **Product Owner** | | | | |
| **Lead Accountant**<br>(Madagascar compliance) | | | | |
| **HR Director**<br>(Workflow validation) | | | | |
| **CTO**<br>(Technical feasibility) | | | | |
| **Legal Counsel**<br>(Regulatory review) | | | | |

---

**Document Control**  
- **Next Review Date:** March 4, 2026  
- **Retention Period:** 10 years (per Madagascar legal requirement)  
- **Distribution:** Product Team, Development Team, Compliance Officer  

---

> ✅ **This PRD incorporates all client requirements:**  
> - Correct CNAPS (pension) vs. OSTIE/OSIE (health) vs. IRSA (tax) distinctions  
> - Minimum wage = 200,000 MGA → CNAPS/OSTIE ceiling = 1,600,000 MGA  
> - Four distinct roles with strict RBAC (Admin, HR, Accountant, Employee)  
> - **Immutable audit trail** for all user actions (Admin oversight)  
> - OHADA-compliant accounting with separation of duties (HR cannot record GL entries)  
> - Full Madagascar regulatory compliance (CNAPS, OSTIE, IRSA, labor code)  
> - Technology stack: Golang (Gin) + PostgreSQL + React + Tailwind  

**Ready for development kickoff upon client approval.**