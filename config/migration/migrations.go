package migration

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"bbscout/config/database"
	"bbscout/models"
)

func InitializeMigrations() {
	database.NewDatabaseConnection()
	db := database.GetDB()

	errMigrationTables := db.AutoMigrate(
		&models.RoleModel{},
		&models.UserModel{},
		&models.FileModel{},
	)
	if errMigrationTables != nil {
		log.Fatalf("failed to migrate tables: %v", errMigrationTables)
	}

	seedRoles(db)

	fmt.Println("Finished migration tables")

}

func seedRoles(db *gorm.DB) {
	fmt.Println("Startting roles seeding")
	roles := []string{"USER"}
	for _, roleName := range roles {
		// Use FirstOrCreate to avoid duplication
		db.FirstOrCreate(&models.RoleModel{}, models.RoleModel{Name: roleName})
	}
}
