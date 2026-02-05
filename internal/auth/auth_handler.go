package auth

import (
	"fmt"
	"net/http"
	"time"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Handler handles authentication requests
type Handler struct {
	repo *Repo
}

// NewHandler creates a new auth handler
func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

// Register handles user registration (admin only)
func (h *Handler) Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "hr": true, "accountant": true, "employee": true}
	if !validRoles[input.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
		EmployeeID:   input.EmployeeID,
		IsActive:     true,
	}

	if err := h.repo.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		EmployeeID: user.EmployeeID,
		IsActive:   user.IsActive,
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.repo.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if account is locked
	locked, err := h.repo.IsAccountLocked(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check account status"})
		return
	}
	if locked {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is locked. Please try again later."})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		// Increment failed login attempts
		_ = h.repo.IncrementFailedLoginAttempts(c.Request.Context(), user.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Update last login
	if err := h.repo.UpdateLastLogin(c.Request.Context(), user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:                    user.ID,
			Email:                 user.Email,
			Role:                  user.Role,
			EmployeeID:            user.EmployeeID,
			Phone:                 user.Phone,
			Address:               user.Address,
			EmergencyContactName:  user.EmergencyContactName,
			EmergencyContactPhone: user.EmergencyContactPhone,
			IsActive:              user.IsActive,
			LastLogin:             user.LastLogin,
		},
	})
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var input RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse refresh token
	token, err := jwt.ParseWithClaims(input.RefreshToken, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(middleware.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Verify user still exists and is active
	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or inactive"})
		return
	}

	// Generate new access token
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, RefreshTokenResponse{
		AccessToken: accessToken,
	})
}

// GetProfile returns the current user's profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:                    user.ID,
		Email:                 user.Email,
		Role:                  user.Role,
		EmployeeID:            user.EmployeeID,
		Phone:                 user.Phone,
		Address:               user.Address,
		EmergencyContactName:  user.EmergencyContactName,
		EmergencyContactPhone: user.EmergencyContactPhone,
		IsActive:              user.IsActive,
		LastLogin:             user.LastLogin,
	})
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ListUsers returns all users (admin only)
func (h *Handler) ListUsers(c *gin.Context) {
	limit := 50
	offset := 0

	users, total, err := h.repo.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:                    user.ID,
			Email:                 user.Email,
			Role:                  user.Role,
			EmployeeID:            user.EmployeeID,
			Phone:                 user.Phone,
			Address:               user.Address,
			EmergencyContactName:  user.EmergencyContactName,
			EmergencyContactPhone: user.EmergencyContactPhone,
			IsActive:              user.IsActive,
			LastLogin:             user.LastLogin,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  userResponses,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateProfile updates current user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input UpdateProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	user, err := h.repo.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields if provided
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}
	if input.Address != nil {
		user.Address = input.Address
	}
	if input.EmergencyContactName != nil {
		user.EmergencyContactName = input.EmergencyContactName
	}
	if input.EmergencyContactPhone != nil {
		user.EmergencyContactPhone = input.EmergencyContactPhone
	}

	user.UpdatedAt = time.Now()

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:                    user.ID,
		Email:                 user.Email,
		Role:                  user.Role,
		EmployeeID:            user.EmployeeID,
		Phone:                 user.Phone,
		Address:               user.Address,
		EmergencyContactName:  user.EmergencyContactName,
		EmergencyContactPhone: user.EmergencyContactPhone,
		IsActive:              user.IsActive,
		UpdatedAt:             &user.UpdatedAt,
	})
}

// GetPreferences retrieves current user's preferences
func (h *Handler) GetPreferences(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	prefs, err := h.repo.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
		return
	}

	response := PreferencesResponse{}
	response.Notifications.EmailEnabled = prefs.EmailNotifications
	response.Notifications.PushEnabled = prefs.PushNotifications
	response.Notifications.LeaveUpdates = prefs.LeaveUpdates
	response.Notifications.PayrollUpdates = prefs.PayrollUpdates
	response.Notifications.SystemUpdates = prefs.SystemUpdates
	response.Display.Theme = prefs.Theme
	response.Display.Language = prefs.Language
	response.Display.DateFormat = prefs.DateFormat

	c.JSON(http.StatusOK, response)
}

// UpdatePreferences updates current user's preferences
func (h *Handler) UpdatePreferences(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing preferences
	prefs, err := h.repo.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
		return
	}

	// Update fields if provided
	if input.EmailNotifications != nil {
		prefs.EmailNotifications = *input.EmailNotifications
	}
	if input.PushNotifications != nil {
		prefs.PushNotifications = *input.PushNotifications
	}
	if input.LeaveUpdates != nil {
		prefs.LeaveUpdates = *input.LeaveUpdates
	}
	if input.PayrollUpdates != nil {
		prefs.PayrollUpdates = *input.PayrollUpdates
	}
	if input.SystemUpdates != nil {
		prefs.SystemUpdates = *input.SystemUpdates
	}
	if input.Theme != nil {
		prefs.Theme = *input.Theme
	}
	if input.Language != nil {
		prefs.Language = *input.Language
	}
	if input.DateFormat != nil {
		prefs.DateFormat = *input.DateFormat
	}

	if err := h.repo.UpdateUserPreferences(c.Request.Context(), prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	response := PreferencesResponse{}
	response.Notifications.EmailEnabled = prefs.EmailNotifications
	response.Notifications.PushEnabled = prefs.PushNotifications
	response.Notifications.LeaveUpdates = prefs.LeaveUpdates
	response.Notifications.PayrollUpdates = prefs.PayrollUpdates
	response.Notifications.SystemUpdates = prefs.SystemUpdates
	response.Display.Theme = prefs.Theme
	response.Display.Language = prefs.Language
	response.Display.DateFormat = prefs.DateFormat

	c.JSON(http.StatusOK, response)
}
