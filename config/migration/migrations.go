package migration

import (
	"fmt"
	"log"
	"strings"

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
		&models.BillboardTypesModel{},
	)
	if errMigrationTables != nil {
		log.Fatalf("failed to migrate tables: %v", errMigrationTables)
	}

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

	createBillboardSummaryQuery := `
	CREATE OR REPLACE VIEW billboard_summary AS
	SELECT 
	bb.organization_id AS organization_id,
		bb.id AS billboard_id,
		bb.board_code,
		bb.created_by_id,
		bb.location,
		bb.latitude,
		bb.longitude,
		bb.width,
		bb.height,
		bb.unit,
		bb.type,
		bb.price,
		bb.occupied,
        bb.owner,
        bb.owner_contact,
        bb.owner_email,
        bb.closeup_image_id,
        bb.material,
        bb.visibility,
        bb.illumination,
        bb.angel,
        bb.structure,
		bb.parent_board_code,
		bb.city,
		bb.object_type,
		bb.active AS billboard_active,
		bc.id AS campaign_id,
		COALESCE(bc.active, false) AS campaign_active,
		bb.image_id,
		bb.created_at,
		bb.updated_at
	FROM bill_boards bb
	LEFT JOIN billboard_campaign bc 
	ON bb.id = bc.billboard_id
	AND bc.active = true;`

	if err := db.Exec(createBillboardSummaryQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createBillboardWeeklyReportQuery := `
	CREATE OR REPLACE VIEW billboard_upload_day_of_week AS
		SELECT
			o.id AS organization_id,
			o.name AS organization_name,
			o.description AS organization_description,
			o.is_active AS organization_active,
			YEAR(FROM_UNIXTIME(b.created_at)) AS upload_year,
			MONTH(FROM_UNIXTIME(b.created_at)) AS upload_month,
			WEEK(FROM_UNIXTIME(b.created_at), 1) AS upload_week_number,
			DAYNAME(FROM_UNIXTIME(b.created_at)) AS upload_day_name,
			COUNT(*) AS total_uploads
		FROM
			organization o
		LEFT JOIN
			bill_boards b ON b.organization_id = o.id
		GROUP BY
			o.id,
			o.name,
			o.description,
			o.is_active,
			YEAR(FROM_UNIXTIME(b.created_at)),
			MONTH(FROM_UNIXTIME(b.created_at)),
			WEEK(FROM_UNIXTIME(b.created_at), 1),
			DAYNAME(FROM_UNIXTIME(b.created_at));`

	if err := db.Exec(createBillboardWeeklyReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createBillboardMonthlyReportQuery := `
	CREATE OR REPLACE VIEW billboard_upload_monthly_report AS
	SELECT
		o.id AS organization_id,
		o.name AS organization_name,
		o.description AS organization_description,
		o.is_active AS organization_active,
		YEAR(FROM_UNIXTIME(b.created_at)) AS upload_year,
		MONTH(FROM_UNIXTIME(b.created_at)) AS upload_month,
		DATE_FORMAT(FROM_UNIXTIME(b.created_at), '%M') AS upload_month_name,
		COUNT(*) AS total_uploads
	FROM
		organization o
	LEFT JOIN
		bill_boards b ON b.organization_id = o.id
	GROUP BY
		o.id,
		o.name,
		o.description,
		o.is_active,
		YEAR(FROM_UNIXTIME(b.created_at)),
		MONTH(FROM_UNIXTIME(b.created_at)),
		DATE_FORMAT(FROM_UNIXTIME(b.created_at), '%M');`

	if err := db.Exec(createBillboardMonthlyReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createBillboardOrganizationReportQuery := `
	CREATE OR REPLACE VIEW billboard_organization_report AS
		SELECT
			o.id AS organization_id,
			o.name AS organization_name,
			o.description AS organization_description,
			COUNT(DISTINCT b.id) AS total_uploads,
			COUNT(DISTINCT CASE WHEN c.active = TRUE THEN b.id END) AS total_occupied,
			COUNT(DISTINCT CASE WHEN c.active = FALSE OR c.id IS NULL THEN b.id END) AS total_not_occupied,
			COUNT(DISTINCT CASE WHEN DATE(FROM_UNIXTIME(b.created_at)) = CURDATE() THEN b.id END) AS uploaded_today,
			COUNT(DISTINCT CASE WHEN YEAR(FROM_UNIXTIME(b.created_at)) = YEAR(CURDATE())
								AND MONTH(FROM_UNIXTIME(b.created_at)) = MONTH(CURDATE())
								THEN b.id END) AS uploaded_this_month,
			COUNT(DISTINCT CASE WHEN YEAR(FROM_UNIXTIME(b.created_at)) = YEAR(CURDATE()) THEN b.id END) AS uploaded_this_year
		FROM
			bill_boards b
		LEFT JOIN
			billboard_campaign c ON c.billboard_id = b.id
		LEFT JOIN
			organization o ON b.organization_id = o.id
		GROUP BY
			o.id, o.name, o.description;`

	if err := db.Exec(createBillboardOrganizationReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createBillboardOrganizationLocationReportQuery := `
	CREATE OR REPLACE  VIEW billboard_organization_location_report AS
		SELECT
			o.id AS organization_id,
			o.name AS organization_name,
			b.location,
			COUNT(*) AS count_per_location
		FROM
			bill_boards b
		LEFT JOIN
			organization o ON b.organization_id = o.id
		GROUP BY
			o.id, o.name,  b.location;`

	if err := db.Exec(createBillboardOrganizationLocationReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createBillboardOrganizationTypeReportQuery := `
	CREATE OR REPLACE VIEW billboard_organization_type_report AS
		SELECT
			o.id AS organization_id,
			o.name AS organization_name,
			b.type,
			COUNT(*) AS type_count
		FROM
			bill_boards b
		LEFT JOIN
			organization o ON b.organization_id = o.id
		GROUP BY
			o.id, o.name, b.type;`

	if err := db.Exec(createBillboardOrganizationTypeReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createUserBillboardUploadReportReportQuery := `
	CREATE OR REPLACE VIEW  user_billboard_uploads_report AS
		SELECT 
			u.id AS user_id,
			CONCAT_WS(' ', u.first_name, COALESCE(u.middle_name, ''), u.last_name) AS user_name,
			u.email,
			o.id AS organization_id,
			o.name AS organization_name,
			COUNT(b.id) AS billboard_count
		FROM users u
		JOIN bill_boards b ON u.id = b.created_by_id
		JOIN organization o ON b.organization_id = o.id
		GROUP BY u.id, u.first_name, u.middle_name, u.last_name, u.email, o.id, o.name;`

	if err := db.Exec(createUserBillboardUploadReportReportQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createUserBillboardUploadReportWeQuery := `
	CREATE OR REPLACE VIEW user_billboard_uploads_report_week AS
		SELECT 
			u.id AS user_id,
			CONCAT_WS(' ', u.first_name, COALESCE(u.middle_name, ''), u.last_name) AS user_name,
			u.email,
			o.id AS organization_id,
			o.name AS organization_name,
			COUNT(b.id) AS billboard_count,
			YEAR(FROM_UNIXTIME(b.created_at)) AS upload_year,
			MONTH(FROM_UNIXTIME(b.created_at)) AS upload_month,
			WEEK(FROM_UNIXTIME(b.created_at), 1) AS week_number,
			DAYNAME(FROM_UNIXTIME(b.created_at)) AS day_of_week
		FROM users u
		JOIN bill_boards b ON u.id = b.created_by_id
		JOIN organization o ON b.organization_id = o.id
		GROUP BY u.id, u.first_name, u.middle_name, u.last_name, u.email, o.id, o.name, 
				upload_year, upload_month, week_number, day_of_week;`

	if err := db.Exec(createUserBillboardUploadReportWeQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createUserBillboardUploadReportMonthlyQuery := `
	CREATE OR REPLACE VIEW user_billboard_uploads_by_month_year AS
		SELECT 
			u.id AS user_id,
			CONCAT_WS(' ', u.first_name, COALESCE(u.middle_name, ''), u.last_name) AS user_name,
			u.email,
			o.id AS organization_id,
			o.name AS organization_name,
			YEAR(FROM_UNIXTIME(b.created_at)) AS upload_year,
			MONTH(FROM_UNIXTIME(b.created_at)) AS upload_month,
			COUNT(b.id) AS billboard_count
		FROM users u
		JOIN bill_boards b ON u.id = b.created_by_id
		JOIN organization o ON b.organization_id = o.id
		GROUP BY u.id, u.first_name, u.middle_name, u.last_name, u.email, 
				o.id, o.name, upload_year, upload_month;`

	if err := db.Exec(createUserBillboardUploadReportMonthlyQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	createFilesSummaryQuery := `
	CREATE OR REPLACE VIEW files_summary AS
	SELECT 
    fs.id AS file_id,
    fs.file_name ,
    fs.file_extension,
    fs.file_type,
    fs.file_url,
    fs.organization_id,
    fs.file_size,
    fs.created_at,
    fs.updated_at,
    bb.id AS billboard_id
    FROM files fs 
    LEFT JOIN bill_boards bb ON fs.id = bb.closeup_image_id OR fs.id = bb.image_id;`

	if err := db.Exec(createFilesSummaryQuery).Error; err != nil {
		log.Fatalf("failed to create view: %v", err)
	}

	fmt.Println("Finished migration tables")

}
