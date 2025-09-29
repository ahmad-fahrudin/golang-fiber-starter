package database

import (
	"app/src/config"
	"app/src/database/seeders"
	"app/src/utils"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, config.DBUser, config.DBPassword, dbName, config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", errDB)
	}

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func ConnectForSeeder(dbHost, dbName string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, config.DBUser, config.DBPassword, dbName, config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent), // Silent mode for seeders
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database for seeder: %+v", err)
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to connect to database for seeder: %+v", errDB)
	}

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}

// SeederInterface defines the interface for all seeders
type SeederInterface interface {
	Run(db *gorm.DB) error
}

// SeederConfig holds configuration for seeder execution
type SeederConfig struct {
	DB *gorm.DB
}

// NewSeederConfig creates a new seeder configuration
func NewSeederConfig() *SeederConfig {
	dbHost := config.DBHost
	if config.IsProd {
		dbHost = "localhost" // Use localhost in production or adjust as needed
	}

	db := ConnectForSeeder(dbHost, config.DBName)
	return &SeederConfig{
		DB: db,
	}
}

// getAllSeeders returns all available seeders
func (sc *SeederConfig) getAllSeeders() []SeederInterface {
	return []SeederInterface{
		&seeders.UserSeeder{},
		// Add more seeders here as you create them
		// &seeders.ProductSeeder{},
		// &seeders.CategorySeeder{},
	}
}

// RunAllSeeders executes all seeders
func (sc *SeederConfig) RunAllSeeders() error {
	utils.Log.Info("Starting to run all seeders...")

	allSeeders := sc.getAllSeeders()

	for _, seeder := range allSeeders {
		seederName := sc.getSeederName(seeder)
		utils.Log.Infof("Running seeder: %s", seederName)

		if err := seeder.Run(sc.DB); err != nil {
			utils.Log.Errorf("Failed to run seeder %s: %v", seederName, err)
			return fmt.Errorf("failed to run seeder %s: %w", seederName, err)
		}

		utils.Log.Infof("Seeder %s completed successfully", seederName)
	}

	utils.Log.Info("All seeders completed successfully!")
	return nil
}

// RunSpecificSeeder executes a specific seeder by name
func (sc *SeederConfig) RunSpecificSeeder(seederName string) error {
	utils.Log.Infof("Starting to run specific seeder: %s", seederName)

	allSeeders := sc.getAllSeeders()

	for _, seeder := range allSeeders {
		currentSeederName := sc.getSeederName(seeder)
		if strings.EqualFold(currentSeederName, seederName) {
			utils.Log.Infof("Running seeder: %s", currentSeederName)

			if err := seeder.Run(sc.DB); err != nil {
				utils.Log.Errorf("Failed to run seeder %s: %v", currentSeederName, err)
				return fmt.Errorf("failed to run seeder %s: %w", currentSeederName, err)
			}

			utils.Log.Infof("Seeder %s completed successfully", currentSeederName)
			return nil
		}
	}

	return fmt.Errorf("seeder '%s' not found", seederName)
}

// RunMultipleSeeders executes multiple specific seeders
func (sc *SeederConfig) RunMultipleSeeders(seederNames []string) error {
	utils.Log.Infof("Starting to run multiple seeders: %v", seederNames)

	for _, seederName := range seederNames {
		if err := sc.RunSpecificSeeder(seederName); err != nil {
			return err
		}
	}

	utils.Log.Info("All specified seeders completed successfully!")
	return nil
}

// ListAvailableSeeders shows all available seeders
func (sc *SeederConfig) ListAvailableSeeders() {
	utils.Log.Info("Available seeders:")

	allSeeders := sc.getAllSeeders()
	for i, seeder := range allSeeders {
		seederName := sc.getSeederName(seeder)
		utils.Log.Infof("%d. %s", i+1, seederName)
	}
}

// getSeederName extracts seeder name from the struct
func (sc *SeederConfig) getSeederName(seeder SeederInterface) string {
	seederType := reflect.TypeOf(seeder)
	if seederType.Kind() == reflect.Ptr {
		seederType = seederType.Elem()
	}

	// Get the function name and extract seeder name
	fullName := runtime.FuncForPC(reflect.ValueOf(seeder).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return strings.Replace(seederType.Name(), "Seeder", "", 1)
	}

	return seederType.Name()
}

// TruncateTable truncates a table (useful for re-seeding)
func (sc *SeederConfig) TruncateTable(tableName string) error {
	utils.Log.Infof("Truncating table: %s", tableName)

	// Disable foreign key checks temporarily
	if err := sc.DB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		// For PostgreSQL, use different syntax
		if err := sc.DB.Exec("SET session_replication_role = replica").Error; err != nil {
			utils.Log.Warnf("Could not disable foreign key checks: %v", err)
		}
	}

	// Truncate the table
	if err := sc.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName)).Error; err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
	}

	// Re-enable foreign key checks
	if err := sc.DB.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		// For PostgreSQL
		if err := sc.DB.Exec("SET session_replication_role = DEFAULT").Error; err != nil {
			utils.Log.Warnf("Could not re-enable foreign key checks: %v", err)
		}
	}

	utils.Log.Infof("Table %s truncated successfully", tableName)
	return nil
}

// RefreshSeeder truncates table and runs seeder
func (sc *SeederConfig) RefreshSeeder(seederName string, tableName string) error {
	utils.Log.Infof("Refreshing seeder: %s (table: %s)", seederName, tableName)

	// Truncate table first
	if err := sc.TruncateTable(tableName); err != nil {
		return err
	}

	// Run seeder
	return sc.RunSpecificSeeder(seederName)
}
