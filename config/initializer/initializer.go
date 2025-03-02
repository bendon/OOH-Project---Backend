package initializer

import (
	"fmt"

	"github.com/google/uuid"

	"bbscout/models"
	"bbscout/repository"
	"bbscout/utils"
)

func InitializerOperationAccount() {
	fmt.Println("Initializing operation account...")
	orgRepo := repository.NewOrganizationRepository()
	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	userRoleRepo := repository.NewUserRoleRepository()
	userAccountRepo := repository.NewUserAccountRepository()
	permissionRepo := repository.NewPermissionRepository()
	accountPermRepo := repository.NewUserAccountPermissionRepository()
	orgUserRepo := repository.NewOrganizationUserRepository()

	CreatePermissions()

	// confirm whether their is operation account
	found, SearchErr := orgRepo.GetOperationOrganization()

	if SearchErr != nil {
		fmt.Println("Failed to create admin account", SearchErr)
		return
	}

	if found != nil {
		fmt.Printf("Organization Found %v \n", found.Name)
		return
	}

	phone := 25499991299
	country := "KE"
	admin := &models.UserModel{}
	admin.Active = true
	admin.Email = "operations@bbscout.com"
	admin.Gender = 1
	admin.Country = &country
	admin.FirstName = "BBScout"
	admin.LastName = "Operations"
	admin.Phone = &phone
	admin.Verified = true
	admin.Password = utils.HashPassword([]byte("bbscout-password"))

	user, userErr := userRepo.CreateUser(admin)
	if userErr != nil {
		fmt.Println("Failed to create admin account")
		return
	}

	role, roleErr := roleRepo.GetRoleByName("OPERATOR")
	if roleErr != nil {
		fmt.Println("Failed to create admin account")
		return
	}

	if role == nil {
		fmt.Println("Role not found")
		return
	}

	//establish operation organization
	organization := &models.OrganizationModel{}
	organization.AdminId = user.ID
	organization.IsActive = true
	organization.IsOperation = true
	organization.Name = "BBscout Hunter"
	org, orgErr := orgRepo.CreateOrganization(organization)

	if orgErr != nil {
		fmt.Println("Failed to create organization account")
		return
	}

	SeedRoles(org.ID)

	// user role
	userRole := &models.UserRoleModel{}
	userRole.RoleId = role.ID
	userRole.UserId = user.ID
	userRole.OrganizationId = org.ID
	userRoleRepo.CreateUserRole(userRole)

	userAccount := models.UserAccountModel{
		Active:         true,
		OrganizationId: &org.ID,
		UserId:         user.ID,
	}
	fir, firstError := userAccountRepo.CreateUserAccount(&userAccount)
	if firstError != nil {
		fmt.Println("Error creating organization license:", firstError)
		return
	}
	fmt.Printf("org id is %v  for %v \n", fir.OrganizationId, fir.ID)

	//attach organization account permission
	accountPermission, _ := permissionRepo.GetPermissionsByAccount("ORGANIZATION")

	for _, perm := range accountPermission {
		userAccountPermission := models.UserAccountPermissionModel{
			AccountId:      fir.ID,
			OrganizationId: org.ID,
			UserId:         fir.UserId,
			PermissionId:   perm.ID,
		}
		accountPermRepo.CreateUserAccountPermission(&userAccountPermission)
	}

	//Organization User
	orgUSer := &models.OrganizationUserModel{}
	orgUSer.OrganizationId = org.ID
	orgUSer.UserId = user.ID
	orgUserRepo.CreateOrganizationUser(orgUSer)

	fmt.Println("Finished initializing operation account")
}

func CreatePermissions() {
	permRepo := repository.NewPermissionRepository()
	var permissionService = []string{
		"Add Organization",
		"Manage Organization",
		"Delete Organization",
		"Lock Organization",
	}

	for _, item := range permissionService {
		perm := models.PermissionModel{
			Name:    item,
			Type:    "OPERATION",
			Account: "ORGANIZATION",
		}
		permRepo.CreatePermission(&perm)
	}

	var permissionOrganization = []string{
		"Dashboard",
		"View Storage Files",
		"Manage Employees",
		"Add Employees",
		"Staff",
		"Manage Roles",
		"Manage Permissions",
		"Add Roles",
		"Update Roles",
	}

	for _, item := range permissionOrganization {
		perm := models.PermissionModel{
			Name:    item,
			Type:    "ORGANIZATION",
			Account: "ORGANIZATION",
		}
		permRepo.CreatePermission(&perm)
	}

}

func SeedRoles(organizationId uuid.UUID) {
	roleRepo := repository.NewRoleRepository()
	fmt.Println("Startting roles seeding")
	roles := []string{"USER", "ADMIN", "OPERATOR"}
	for _, roleName := range roles {
		// Use FirstOrCreate to avoid duplication
		roleRepo.CreateRole(&models.RoleModel{Name: roleName, OrganizationId: &organizationId})
	}
}
