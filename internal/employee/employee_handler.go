package employee

import (
	"net/http"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles employee requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new employee handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Create handles employee creation (HR/Admin only)
func (h *Handler) Create(c *gin.Context) {
	var input CreateEmployeeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee := &Employee{
		CompanyID:             input.CompanyID,
		FirstName:             input.FirstName,
		LastName:              input.LastName,
		DateOfBirth:           input.DateOfBirth,
		Gender:                input.Gender,
		Nationality:           input.Nationality,
		NationalID:            input.NationalID,
		Position:              input.Position,
		Department:            input.Department,
		HireDate:              input.HireDate,
		ContractType:          input.ContractType,
		GrossSalary:           input.GrossSalary,
		Status:                "active",
		Address:               input.Address,
		Phone:                 input.Phone,
		EmergencyContactName:  input.EmergencyContactName,
		EmergencyContactPhone: input.EmergencyContactPhone,
		ManagerID:             input.ManagerID,
	}

	if err := h.repo.Create(c.Request.Context(), employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// GetByID retrieves an employee by ID
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Check permissions: employees can only view their own profile
	userRole, _ := middleware.GetUserRole(c)
	if userRole == "employee" {
		userID, _ := middleware.GetUserID(c)
		// TODO: Check if userID matches employee's user_id
		_ = userID
	}

	c.JSON(http.StatusOK, employee)
}

// Update updates an employee (HR/Admin only)
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var input UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Update fields if provided
	if input.FirstName != nil {
		employee.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		employee.LastName = *input.LastName
	}
	if input.DateOfBirth != nil {
		employee.DateOfBirth = input.DateOfBirth
	}
	if input.Gender != nil {
		employee.Gender = *input.Gender
	}
	if input.Nationality != nil {
		employee.Nationality = *input.Nationality
	}
	if input.NationalID != nil {
		employee.NationalID = *input.NationalID
	}
	if input.Position != nil {
		employee.Position = *input.Position
	}
	if input.Department != nil {
		employee.Department = *input.Department
	}
	if input.ContractType != nil {
		employee.ContractType = *input.ContractType
	}
	if input.GrossSalary != nil {
		employee.GrossSalary = *input.GrossSalary
	}
	if input.Status != nil {
		employee.Status = *input.Status
	}
	if input.Address != nil {
		employee.Address = *input.Address
	}
	if input.Phone != nil {
		employee.Phone = *input.Phone
	}
	if input.EmergencyContactName != nil {
		employee.EmergencyContactName = *input.EmergencyContactName
	}
	if input.EmergencyContactPhone != nil {
		employee.EmergencyContactPhone = *input.EmergencyContactPhone
	}
	if input.ManagerID != nil {
		employee.ManagerID = input.ManagerID
	}

	if err := h.repo.Update(c.Request.Context(), employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// Delete soft deletes an employee (Admin only)
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// List retrieves employees with filtering and pagination
func (h *Handler) List(c *gin.Context) {
	var query EmployeeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employees, total, err := h.repo.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list employees"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"total":     total,
		"limit":     query.Limit,
		"offset":    query.Offset,
	})
}

// GetDepartments retrieves all unique departments
func (h *Handler) GetDepartments(c *gin.Context) {
	departments, err := h.repo.GetDepartments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get departments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"departments": departments})
}

// GetPositions retrieves all unique positions
func (h *Handler) GetPositions(c *gin.Context) {
	positions, err := h.repo.GetPositions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get positions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

// GetSubordinates retrieves all employees reporting to a manager
func (h *Handler) GetSubordinates(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
		return
	}

	subordinates, err := h.repo.GetSubordinates(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subordinates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subordinates": subordinates})
}
