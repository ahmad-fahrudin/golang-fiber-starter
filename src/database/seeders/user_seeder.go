package seeders

import (
	"app/src/model"
	"app/src/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder handles user data seeding
type UserSeeder struct{}

// Run executes the user seeder
func (s *UserSeeder) Run(db *gorm.DB) error {
	utils.Log.Info("Seeding users...")

	users := []model.User{
		{
			Name:          "Admin User",
			Email:         "admin@gmail.com",
			Password:      s.hashPassword("passsword"),
			Role:          "admin",
			VerifiedEmail: true,
		},
	}

	createdCount := 0
	skippedCount := 0

	for _, user := range users {
		// Check if user already exists
		var existingUser model.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// User doesn't exist, create it
				if err := db.Create(&user).Error; err != nil {
					utils.Log.Errorf("Failed to create user %s: %v", user.Email, err)
					return err
				}
				createdCount++
			} else {
				utils.Log.Errorf("Error checking user existence %s: %v", user.Email, err)
				return err
			}
		} else {
			skippedCount++
		}
	}

	utils.Log.Infof("User seeder completed: %d created, %d skipped", createdCount, skippedCount)
	return nil
}

// hashPassword hashes a password using bcrypt
func (s *UserSeeder) hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Errorf("Failed to hash password: %v", err)
		return password // Return original password if hashing fails (not recommended for production)
	}
	return string(hashedPassword)
}
