# Missing API Endpoints Request

This document lists all the API endpoints that are currently not available in the backend but are needed by the frontend application. These endpoints should be created to fully integrate the real APIs and remove mock data dependencies.

## Date Created
February 5, 2026

## Status
🔴 **CRITICAL** - These endpoints are required for the application to function properly with real data.

---

## 1. Dashboard Statistics Endpoint

**Endpoint:** `GET /api/v1/dashboard/stats`

**Purpose:** Get aggregated dashboard statistics for the current user based on their role.

**Access:** All authenticated users

**Response:**
```json
{
  "totalEmployees": 8,
  "onLeave": 1,
  "newRequests": 3,
  "attendanceRate": 87.5,
  "pendingLeaves": 2,
  "monthlyPayroll": 24500000
}
```

**Why Needed:**
- Currently, the frontend has to make multiple API calls and calculate these stats client-side
- This is inefficient and can lead to performance issues
- The dashboard needs this data on every page load

**Alternative:**
If this endpoint cannot be created, the frontend will continue to calculate these from existing endpoints:
- `GET /employees` for totalEmployees
- `GET /employees?status=on_leave` for onLeave
- `GET /leaves?status=pending` for pendingLeaves and newRequests
- `GET /attendance` + calculation for attendanceRate
- Sum of all employee salaries for monthlyPayroll

---

## 2. Support Ticket Submission (Already Exists ✅)

**Endpoint:** `POST /api/v1/support`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists and is documented in the API docs. No action needed.

---

## 3. User Management Endpoints (Partially Missing)

### GET /api/v1/auth/users
**Status:** ✅ **EXISTS** in API documentation

### Additional User Management Endpoints Needed:

**Endpoint:** `PUT /api/v1/auth/users/:id`

**Purpose:** Update user information (role, status, etc.)

**Access:** Admin only

**Request Body:**
```json
{
  "role": "hr",
  "is_active": true
}
```

**Endpoint:** `DELETE /api/v1/auth/users/:id`

**Purpose:** Delete/deactivate a user account

**Access:** Admin only

---

## 4. Subordinates Hierarchy Endpoint (Already Exists ✅)

**Endpoint:** `GET /api/v1/employees/:id/subordinates`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists but is not currently used in the frontend. It should be integrated for manager views.

---

## 5. Attendance Statistics Summary

**Endpoint:** `GET /api/v1/attendance/summary`

**Purpose:** Get a daily or weekly summary of attendance across the organization

**Access:** HR, Admin

**Query Parameters:**
- `date` - Specific date (default: today)
- `period` - 'day' | 'week' | 'month'

**Response:**
```json
{
  "date": "2026-02-05",
  "totalEmployees": 8,
  "present": 5,
  "late": 1,
  "absent": 2,
  "overtime": 0,
  "attendanceRate": 75.0
}
```

**Why Needed:**
- The attendance page shows daily statistics but has to calculate them client-side
- This endpoint would provide optimized server-side calculation

**Alternative:**
Continue using `GET /attendance?start_date=X&end_date=X` and calculate client-side.

---

## 6. Leave Balance for Multiple Employees

**Endpoint:** `GET /api/v1/leaves/balances`

**Purpose:** Get leave balances for all employees or filtered list

**Access:** HR, Admin

**Query Parameters:**
- `year` - Year for balance calculation (default: current year)
- `department` - Filter by department
- `employee_ids` - Comma-separated list of employee IDs

**Response:**
```json
[
  {
    "employee_id": "emp-001",
    "year": 2024,
    "annual_entitlement": 30,
    "annual_used": 3,
    "annual_remaining": 27,
    "sick_used": 2
  },
  ...
]
```

**Why Needed:**
- HR needs to see leave balances for multiple employees at once
- Currently requires individual API calls for each employee

**Alternative:**
Make individual calls to `GET /leaves/balance/:employee_id` for each employee (inefficient).

---

## 7. Payroll Reconciliation Report (Already Exists ✅)

**Endpoint:** `GET /api/v1/payroll/reconciliation`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists and should be integrated into the payroll page.

---

## 8. Declaration Form Generation (Partially Exists ✅)

**Endpoints:**
- `GET /api/v1/declarations/cnaps/generate?month=YYYY-MM` ✅ **EXISTS**
- `GET /api/v1/declarations/ostie/generate?month=YYYY-MM` ✅ **EXISTS**
- `GET /api/v1/declarations/irsa/generate?month=YYYY-MM` ✅ **EXISTS**

These endpoints exist and should be integrated into the declarations page for automatic form generation.

---

## 9. Declaration Form Populate (Already Exists ✅)

**Endpoint:** `POST /api/v1/declarations/:id/populate`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists to populate declaration data from approved payrolls.

---

## 10. Performance Reports (Already Exists ✅)

**Endpoint:** `GET /api/v1/kpi/reports/:employee_id`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists and should be integrated for comprehensive performance reports.

---

## 11. Audit Log Export (Already Exists ✅)

**Endpoint:** `GET /api/v1/audit/logs/export`

**Status:** ✅ **EXISTS** in API documentation

This endpoint exists and should be integrated for CSV export functionality.

---

## 12. Bulk Operations (Nice to Have)

### Bulk Clock-In/Clock-Out
**Endpoint:** `POST /api/v1/attendance/bulk-clock-in`

**Purpose:** Allow HR to clock in multiple employees at once (e.g., for events)

**Priority:** 🟡 Low

### Bulk Leave Approval
**Endpoint:** `POST /api/v1/leaves/bulk-approve`

**Purpose:** Approve multiple leave requests at once

**Priority:** 🟡 Low

---

## Summary

### ✅ Endpoints That Already Exist (should be integrated)
1. Employee subordinates
2. Payroll reconciliation
3. Declaration form generation (CNAPS, OSTIE, IRSA)
4. Declaration populate
5. Performance reports
6. Audit log export
7. Support ticket submission

### 🔴 Critical Missing Endpoints
1. **Dashboard statistics** - High priority
2. **Attendance summary** - Medium priority
3. **Leave balance bulk query** - Medium priority

### 🟡 Nice-to-Have Endpoints
1. User update/delete endpoints
2. Bulk operations for attendance and leaves

---

## Recommendations

1. **Immediate Action:** Create the Dashboard Statistics endpoint as it's called on every dashboard load
2. **Short Term:** Implement Attendance Summary and Leave Balance bulk query for better performance
3. **Medium Term:** Add user management endpoints for complete admin functionality
4. **Long Term:** Consider bulk operations for improved UX

---

## Notes for Backend Development

- All new endpoints should follow the existing API patterns (authentication, error handling, etc.)
- Consider pagination for bulk queries
- Add appropriate role-based access control
- Document all new endpoints in the API_DOCUMENTATION.md file
- Follow Madagascar-specific business logic where applicable (salary constraints, leave calculations, etc.)

---

**Document Maintained By:** Frontend Development Team  
**Last Updated:** February 5, 2026  
**Next Review:** When backend implements missing endpoints
