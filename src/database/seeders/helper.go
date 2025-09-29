package seeders

import (
	"app/src/utils"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// SeederHelper provides common utilities for seeders
type SeederHelper struct{}

// NewSeederHelper creates a new seeder helper
func NewSeederHelper() *SeederHelper {
	return &SeederHelper{}
}

// CheckRecordExists checks if a record exists based on a condition
func (h *SeederHelper) CheckRecordExists(db *gorm.DB, model interface{}, condition string, args interface{}) bool {
	err := db.Where(condition, args).First(model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		utils.Log.Errorf("Error checking record existence: %v", err)
		return true // Return true to be safe and skip creation
	}
	return true
}

// CreateRecordIfNotExists creates a record if it doesn't exist
func (h *SeederHelper) CreateRecordIfNotExists(db *gorm.DB, record interface{}, checkCondition string, args interface{}, identifier string) error {
	// Check if record exists
	err := db.Where(checkCondition, args).First(record).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record doesn't exist, create it
			if err := db.Create(record).Error; err != nil {
				utils.Log.Errorf("Failed to create record %s: %v", identifier, err)
				return err
			}
			utils.Log.Infof("Created record: %s", identifier)
		} else {
			utils.Log.Errorf("Error checking record existence %s: %v", identifier, err)
			return err
		}
	} else {
		utils.Log.Infof("Record %s already exists, skipping...", identifier)
	}

	return nil
}

// GenerateRandomString generates a random string of specified length
func (h *SeederHelper) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomEmail generates a random email
func (h *SeederHelper) GenerateRandomEmail() string {
	domains := []string{"example.com", "test.com", "sample.org", "demo.net"}
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	username := h.GenerateRandomString(8)
	domain := domains[seededRand.Intn(len(domains))]

	return username + "@" + domain
}

// GetRandomItem returns a random item from a slice of strings
func (h *SeederHelper) GetRandomItem(items []string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return items[seededRand.Intn(len(items))]
}

// TruncateAndResetTable truncates a table and resets its auto-increment
func (h *SeederHelper) TruncateAndResetTable(db *gorm.DB, tableName string) error {
	utils.Log.Infof("Truncating and resetting table: %s", tableName)

	// For PostgreSQL
	if err := db.Exec("TRUNCATE TABLE " + tableName + " RESTART IDENTITY CASCADE").Error; err != nil {
		// Fallback for other databases
		if err := db.Exec("DELETE FROM " + tableName).Error; err != nil {
			utils.Log.Errorf("Failed to truncate table %s: %v", tableName, err)
			return err
		}
	}

	utils.Log.Infof("Table %s truncated and reset successfully", tableName)
	return nil
}

// GetRandomInt generates a random integer between min and max
func (h *SeederHelper) GetRandomInt(min, max int) int {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(max-min+1) + min
}

// GetRandomBool generates a random boolean value
func (h *SeederHelper) GetRandomBool() bool {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(2) == 1
}

// GetRandomDate generates a random date within the last specified days
func (h *SeederHelper) GetRandomDate(daysBack int) time.Time {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	days := seededRand.Intn(daysBack)
	return time.Now().AddDate(0, 0, -days)
}
