package migration

import (
	"fmt"
	"log"
	"strings"

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
		&models.OrganizationModel{},
		&models.UserAccountModel{},
		&models.PermissionModel{},
		&models.UserAccountPermissionModel{},
		&models.UserRoleModel{},
		&models.OrganizationUserModel{},
		&models.BillboardModel{},
		&models.BillboardSequenceModel{},
		&models.BillboardCampaignModel{},
	)
	if errMigrationTables != nil {
		log.Fatalf("failed to migrate tables: %v", errMigrationTables)
	}

	seedRoles(db)

	createOrganizationUserSummaryQuery := `
	CREATE OR REPLACE VIEW organization_user_summary AS SELECT 
			ou.id AS organization_user_id,
			ou.organization_id,
			o.name AS organization_name,
			u.id AS user_id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.email,
			u.phone,
			u.country,
			u.gender,
			u.verified,
			u.active,
			r.id AS role_id,
			r.name AS role_name,
			ou.created_at,
			ou.updated_at
		FROM organization_user ou
		JOIN users u ON ou.user_id = u.id
		JOIN organization o ON ou.organization_id = o.id
		LEFT JOIN user_role ur ON ur.user_id = u.id AND ur.organization_id = ou.organization_id
		LEFT JOIN roles r ON ur.role_id = r.id;`

	if err := db.Exec(createOrganizationUserSummaryQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createOrganizationAnalyticsQuery := `
	CREATE OR REPLACE VIEW organization_user_analytics AS
		SELECT 
			o.id AS organization_id,
			o.name AS organization_name,
			COUNT(DISTINCT ou.user_id) AS no_of_users,
			COUNT(DISTINCT CASE WHEN u.gender = 1 THEN u.id END) AS male_count,
			COUNT(DISTINCT CASE WHEN u.gender = 2 THEN u.id END) AS female_count,
			COUNT(DISTINCT CASE WHEN u.gender = 3 THEN u.id END) AS transgender_count,
			COUNT(DISTINCT ur.role_id) AS no_of_roles,
			COUNT(DISTINCT CASE 
				WHEN FROM_UNIXTIME(ou.created_at) >= DATE_FORMAT(NOW(), '%Y-%m-01') 
				THEN ou.user_id 
			END) AS joined_this_month,
			COUNT(DISTINCT CASE 
				WHEN u.verified = 1 THEN u.id 
			END) AS verified_users,
			COUNT(DISTINCT CASE 
				WHEN u.active = 1 THEN u.id 
			END) AS active_users,
			MAX(ou.created_at) AS last_user_joined_at
		FROM organization_user ou
		JOIN users u ON ou.user_id = u.id
		JOIN organization o ON ou.organization_id = o.id
		LEFT JOIN user_role ur ON ur.user_id = u.id AND ur.organization_id = o.id
		GROUP BY o.id, o.name;`

	if err := db.Exec(createOrganizationAnalyticsQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	dropTriggerboardSequencesTrigger := `DROP TRIGGER IF EXISTS before_insert_bill_board_sequences;`
	if err := db.Exec(dropTriggerboardSequencesTrigger).Error; err != nil {
		log.Fatalf("failed to drop trigger: %v", err)
	}

	triggerSQL := `
		
		DELIMITER $$

		CREATE TRIGGER before_insert_bill_board_sequences
		BEFORE INSERT ON bill_board_sequences
		FOR EACH ROW
		BEGIN
			DECLARE max_board_number BIGINT;
			SELECT COALESCE(MAX(board_number), 0) INTO max_board_number
			FROM bill_board_sequences
			WHERE organization_id = NEW.organization_id;

			SET NEW.board_number = max_board_number + 1;
		END$$

		DELIMITER ;
	`

	cleanSQL := strings.ReplaceAll(triggerSQL, "DELIMITER $$", "")
	cleanSQL = strings.ReplaceAll(cleanSQL, "DELIMITER ;", "")
	cleanSQL = strings.ReplaceAll(cleanSQL, "$$", ";") // Replace $$ with ;
	if err := db.Exec(cleanSQL).Error; err != nil {
		log.Fatalf("failed to create trigger: %v", err)
	}

	fmt.Println("Finished migration tables")

}

func seedRoles(db *gorm.DB) {
	fmt.Println("Startting roles seeding")
	roles := []string{"USER", "ADMIN", "OPERATOR"}
	for _, roleName := range roles {
		// Use FirstOrCreate to avoid duplication
		db.FirstOrCreate(&models.RoleModel{}, models.RoleModel{Name: roleName})
	}
}
