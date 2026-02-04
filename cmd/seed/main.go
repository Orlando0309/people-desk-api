package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go-server/internal/auth"
	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/employee"
	"go-server/internal/payroll"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	switch command {
	case "all":
		if err := seedAll(database); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	case "users":
		if err := seedUsers(database); err != nil {
			log.Fatalf("Seeding users failed: %v", err)
		}
	case "employees":
		if err := seedEmployees(database); err != nil {
			log.Fatalf("Seeding employees failed: %v", err)
		}
	case "irsa":
		if err := seedIRSATaxBrackets(database); err != nil {
			log.Fatalf("Seeding IRSA tax brackets failed: %v", err)
		}
	case "payroll-config":
		if err := seedPayrollConfigurations(database); err != nil {
			log.Fatalf("Seeding payroll configurations failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/seed/main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  all             Seed all default data")
	fmt.Println("  users           Seed default users")
	fmt.Println("  employees       Seed default employees")
	fmt.Println("  irsa            Seed default IRSA tax brackets")
	fmt.Println("  payroll-config  Seed default payroll configurations")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/seed/main.go all")
	fmt.Println("  go run cmd/seed/main.go users")
	fmt.Println("  go run cmd/seed/main.go irsa")
}

func seedAll(database *gorm.DB) error {
	fmt.Println("Seeding all default data...")

	if err := seedUsers(database); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := seedEmployees(database); err != nil {
		return fmt.Errorf("failed to seed employees: %w", err)
	}

	if err := seedIRSATaxBrackets(database); err != nil {
		return fmt.Errorf("failed to seed IRSA tax brackets: %w", err)
	}

	if err := seedPayrollConfigurations(database); err != nil {
		return fmt.Errorf("failed to seed payroll configurations: %w", err)
	}

	fmt.Println("Successfully seeded all default data!")
	return nil
}

func seedUsers(database *gorm.DB) error {
	fmt.Println("Seeding default users...")

	// Default admin user
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := auth.User{
		Email:        "admin@peopledesk.com",
		PasswordHash: string(adminPassword),
		Role:         "admin",
		IsActive:     true,
	}

	// Default HR user
	hrPassword, _ := bcrypt.GenerateFromPassword([]byte("hr123"), bcrypt.DefaultCost)
	hrUser := auth.User{
		Email:        "hr@peopledesk.com",
		PasswordHash: string(hrPassword),
		Role:         "hr",
		IsActive:     true,
	}

	// Default Accountant user
	accountantPassword, _ := bcrypt.GenerateFromPassword([]byte("accountant123"), bcrypt.DefaultCost)
	accountantUser := auth.User{
		Email:        "accountant@peopledesk.com",
		PasswordHash: string(accountantPassword),
		Role:         "accountant",
		IsActive:     true,
	}

	// Default Employee user
	employeePassword, _ := bcrypt.GenerateFromPassword([]byte("employee123"), bcrypt.DefaultCost)
	employeeUser := auth.User{
		Email:        "employee@peopledesk.com",
		PasswordHash: string(employeePassword),
		Role:         "employee",
		IsActive:     true,
	}

	users := []auth.User{adminUser, hrUser, accountantUser, employeeUser}

	for _, user := range users {
		var existingUser auth.User
		result := database.Where("email = ?", user.Email).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := database.Create(&user).Error; err != nil {
					return fmt.Errorf("failed to create user %s: %w", user.Email, err)
				}
				fmt.Printf("  Created user: %s (role: %s)\n", user.Email, user.Role)
			} else {
				return fmt.Errorf("failed to query user %s: %w", user.Email, result.Error)
			}
		} else {
			fmt.Printf("  User already exists: %s (role: %s)\n", user.Email, user.Role)
		}
	}

	fmt.Println("Default users seeded successfully!")
	return nil
}

func seedEmployees(database *gorm.DB) error {
	fmt.Println("Seeding default employees...")

	// Default company ID (you may want to make this configurable)
	companyID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	now := time.Now()

	// Default employees
	employees := []employee.Employee{
		{
			CompanyID:             companyID,
			FirstName:             "Admin",
			LastName:              "User",
			Gender:                "other",
			Nationality:           "Malagasy",
			Position:              "System Administrator",
			Department:            "IT",
			HireDate:              now.AddDate(-2, 0, 0),
			ContractType:          "permanent",
			GrossSalary:           1500000.00,
			Status:                "active",
			Address:               "Antananarivo, Madagascar",
			Phone:                 "+261 34 00 000 01",
			EmergencyContactName:  "Emergency Contact 1",
			EmergencyContactPhone: "+261 34 00 000 99",
		},
		{
			CompanyID:             companyID,
			FirstName:             "HR",
			LastName:              "Manager",
			Gender:                "female",
			Nationality:           "Malagasy",
			Position:              "HR Manager",
			Department:            "Human Resources",
			HireDate:              now.AddDate(-1, 6, 0),
			ContractType:          "permanent",
			GrossSalary:           1200000.00,
			Status:                "active",
			Address:               "Antananarivo, Madagascar",
			Phone:                 "+261 34 00 000 02",
			EmergencyContactName:  "Emergency Contact 2",
			EmergencyContactPhone: "+261 34 00 000 98",
		},
		{
			CompanyID:             companyID,
			FirstName:             "Accountant",
			LastName:              "Specialist",
			Gender:                "male",
			Nationality:           "Malagasy",
			Position:              "Senior Accountant",
			Department:            "Finance",
			HireDate:              now.AddDate(-1, 0, 0),
			ContractType:          "permanent",
			GrossSalary:           1100000.00,
			Status:                "active",
			Address:               "Antananarivo, Madagascar",
			Phone:                 "+261 34 00 000 03",
			EmergencyContactName:  "Emergency Contact 3",
			EmergencyContactPhone: "+261 34 00 000 97",
		},
		{
			CompanyID:             companyID,
			FirstName:             "John",
			LastName:              "Doe",
			Gender:                "male",
			Nationality:           "Malagasy",
			Position:              "Software Developer",
			Department:            "IT",
			HireDate:              now.AddDate(-0, 6, 0),
			ContractType:          "permanent",
			GrossSalary:           800000.00,
			Status:                "active",
			Address:               "Antananarivo, Madagascar",
			Phone:                 "+261 34 00 000 04",
			EmergencyContactName:  "Emergency Contact 4",
			EmergencyContactPhone: "+261 34 00 000 96",
		},
		{
			CompanyID:             companyID,
			FirstName:             "Jane",
			LastName:              "Smith",
			Gender:                "female",
			Nationality:           "Malagasy",
			Position:              "Marketing Specialist",
			Department:            "Marketing",
			HireDate:              now.AddDate(-0, 3, 0),
			ContractType:          "permanent",
			GrossSalary:           700000.00,
			Status:                "active",
			Address:               "Antananarivo, Madagascar",
			Phone:                 "+261 34 00 000 05",
			EmergencyContactName:  "Emergency Contact 5",
			EmergencyContactPhone: "+261 34 00 000 95",
		},
	}

	for _, emp := range employees {
		var existingEmployee employee.Employee
		result := database.Where("first_name = ? AND last_name = ?", emp.FirstName, emp.LastName).First(&existingEmployee)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := database.Create(&emp).Error; err != nil {
					return fmt.Errorf("failed to create employee %s %s: %w", emp.FirstName, emp.LastName, err)
				}
				fmt.Printf("  Created employee: %s %s (%s)\n", emp.FirstName, emp.LastName, emp.Position)
			} else {
				return fmt.Errorf("failed to query employee %s %s: %w", emp.FirstName, emp.LastName, result.Error)
			}
		} else {
			fmt.Printf("  Employee already exists: %s %s (%s)\n", emp.FirstName, emp.LastName, emp.Position)
		}
	}

	fmt.Println("Default employees seeded successfully!")
	return nil
}

func seedIRSATaxBrackets(database *gorm.DB) error {
	fmt.Println("Seeding default IRSA tax brackets...")

	// Default system user ID for created_by
	systemUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// IRSA tax brackets for Madagascar (example values - adjust according to actual regulations)
	taxBrackets := []payroll.IRSATaxBracket{
		{
			MinIncome:   0,
			MaxIncome:   float64Ptr(350000),
			TaxRate:     0.00,
			MinTax:      0,
			BracketName: "Tranche 1 - 0%",
			IsActive:    true,
			SortOrder:   1,
			CreatedBy:   systemUserID,
		},
		{
			MinIncome:   350001,
			MaxIncome:   float64Ptr(400000),
			TaxRate:     0.05,
			MinTax:      0,
			BracketName: "Tranche 2 - 5%",
			IsActive:    true,
			SortOrder:   2,
			CreatedBy:   systemUserID,
		},
		{
			MinIncome:   400001,
			MaxIncome:   float64Ptr(500000),
			TaxRate:     0.10,
			MinTax:      2500,
			BracketName: "Tranche 3 - 10%",
			IsActive:    true,
			SortOrder:   3,
			CreatedBy:   systemUserID,
		},
		{
			MinIncome:   500001,
			MaxIncome:   float64Ptr(600000),
			TaxRate:     0.15,
			MinTax:      12500,
			BracketName: "Tranche 4 - 15%",
			IsActive:    true,
			SortOrder:   4,
			CreatedBy:   systemUserID,
		},
		{
			MinIncome:   600001,
			MaxIncome:   nil, // No upper limit
			TaxRate:     0.20,
			MinTax:      27500,
			BracketName: "Tranche 5 - 20%",
			IsActive:    true,
			SortOrder:   5,
			CreatedBy:   systemUserID,
		},
	}

	for _, bracket := range taxBrackets {
		var existingBracket payroll.IRSATaxBracket
		var result *gorm.DB
		if bracket.MaxIncome == nil {
			result = database.Where("min_income = ? AND max_income IS NULL", bracket.MinIncome).First(&existingBracket)
		} else {
			result = database.Where("min_income = ? AND max_income = ?", bracket.MinIncome, *bracket.MaxIncome).First(&existingBracket)
		}

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := database.Create(&bracket).Error; err != nil {
					return fmt.Errorf("failed to create tax bracket: %w", err)
				}
				maxStr := "No upper limit"
				if bracket.MaxIncome != nil {
					maxStr = fmt.Sprintf("%.2f", *bracket.MaxIncome)
				}
				fmt.Printf("  Created tax bracket: %.2f - %s Ar (rate: %.0f%%)\n", bracket.MinIncome, maxStr, bracket.TaxRate*100)
			} else {
				return fmt.Errorf("failed to query tax bracket: %w", result.Error)
			}
		} else {
			maxStr := "No upper limit"
			if bracket.MaxIncome != nil {
				maxStr = fmt.Sprintf("%.2f", *bracket.MaxIncome)
			}
			fmt.Printf("  Tax bracket already exists: %.2f - %s Ar (rate: %.0f%%)\n", bracket.MinIncome, maxStr, bracket.TaxRate*100)
		}
	}

	fmt.Println("Default IRSA tax brackets seeded successfully!")
	return nil
}

func seedPayrollConfigurations(database *gorm.DB) error {
	fmt.Println("Seeding default payroll configurations...")

	// Default system user ID for created_by
	systemUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Default payroll configurations
	configs := []payroll.PayrollConfiguration{
		{
			Key:         "minimum_wage",
			Value:       "200000",
			Description: "Minimum wage (Ar) per month",
			DataType:    "number",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "standard_working_hours",
			Value:       "173.33",
			Description: "Standard working hours per month (40 hours/week)",
			DataType:    "number",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "overtime_multiplier_1",
			Value:       "1.30",
			Description: "Overtime multiplier for first 8 hours (30%)",
			DataType:    "number",
			Category:    "overtime",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "overtime_multiplier_2",
			Value:       "1.50",
			Description: "Overtime multiplier for hours beyond 8 (50%)",
			DataType:    "number",
			Category:    "overtime",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "cnaps_rate_employee",
			Value:       "0.01",
			Description: "CNAPS contribution rate for employee (1%)",
			DataType:    "number",
			Category:    "social_security",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "cnaps_rate_employer",
			Value:       "0.13",
			Description: "CNAPS contribution rate for employer (13%)",
			DataType:    "number",
			Category:    "social_security",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "ostie_rate_employee",
			Value:       "0.01",
			Description: "OSTIE contribution rate for employee (1%)",
			DataType:    "number",
			Category:    "social_security",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "ostie_rate_employer",
			Value:       "0.05",
			Description: "OSTIE contribution rate for employer (5%)",
			DataType:    "number",
			Category:    "social_security",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "medical_allowance",
			Value:       "50000",
			Description: "Medical allowance per month (Ar)",
			DataType:    "number",
			Category:    "allowances",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "transport_allowance",
			Value:       "30000",
			Description: "Transport allowance per month (Ar)",
			DataType:    "number",
			Category:    "allowances",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "housing_allowance",
			Value:       "100000",
			Description: "Housing allowance per month (Ar)",
			DataType:    "number",
			Category:    "allowances",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "family_allowance_rate",
			Value:       "0.05",
			Description: "Family allowance rate (5%)",
			DataType:    "number",
			Category:    "allowances",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "annual_leave_days",
			Value:       "30",
			Description: "Annual leave days per year",
			DataType:    "number",
			Category:    "leave",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "sick_leave_days",
			Value:       "15",
			Description: "Sick leave days per year",
			DataType:    "number",
			Category:    "leave",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "payroll_processing_day",
			Value:       "25",
			Description: "Day of month to process payroll",
			DataType:    "number",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "currency_code",
			Value:       "MGA",
			Description: "Currency code for payroll",
			DataType:    "string",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "tax_year_start",
			Value:       "01-01",
			Description: "Tax year start date (MM-DD)",
			DataType:    "string",
			Category:    "tax",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "enable_auto_payroll",
			Value:       "false",
			Description: "Enable automatic payroll processing",
			DataType:    "boolean",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "irsa_enabled",
			Value:       "true",
			Description: "Enable IRSA tax calculation",
			DataType:    "boolean",
			Category:    "tax",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
		{
			Key:         "rounding_method",
			Value:       "nearest",
			Description: "Rounding method for calculations (nearest, up, down)",
			DataType:    "string",
			Category:    "general",
			IsActive:    true,
			CreatedBy:   systemUserID,
		},
	}

	for _, config := range configs {
		var existingConfig payroll.PayrollConfiguration
		result := database.Where("key = ?", config.Key).First(&existingConfig)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				if err := database.Create(&config).Error; err != nil {
					return fmt.Errorf("failed to create config %s: %w", config.Key, err)
				}
				fmt.Printf("  Created config: %s = %s\n", config.Key, config.Value)
			} else {
				return fmt.Errorf("failed to query config %s: %w", config.Key, result.Error)
			}
		} else {
			fmt.Printf("  Config already exists: %s = %s\n", config.Key, config.Value)
		}
	}

	fmt.Println("Default payroll configurations seeded successfully!")
	return nil
}

func float64Ptr(f float64) *float64 {
	return &f
}
