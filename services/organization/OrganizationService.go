package services

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	emails "bbscout/services/email"
	"bbscout/types"
	"bbscout/utils"
)

func GetOrganizationProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	orgRepo := repository.NewOrganizationRepository()
	org, err := orgRepo.GetOrganizationById(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	// organization profile not found
	if org == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "organization not found")
	}
	return c.Status(fiber.StatusOK).JSON(org)

}

func NewOrganizationStaff(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	var payload types.CreateStaffRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}
	orgRepo := repository.NewOrganizationRepository()
	userRepo := repository.NewUserRepository()
	orgUserRepo := repository.NewOrganizationUserRepository()

	existUser, err := userRepo.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
	}
	if existUser != nil {

		//check if user is already a staff in the organization
		existsOrgUser, err := orgUserRepo.ExistUserOrganizationByUserId(existUser.ID, user.Accessor)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
		}

		// if user is not a staff in the organization
		if !existsOrgUser {
			orgUser := &models.OrganizationUserModel{
				OrganizationId: user.Accessor,
				UserId:         existUser.ID,
				CreatedById:    &user.OwnerID,
			}
			_, err := orgUserRepo.CreateOrganizationUser(orgUser)
			if err != nil {
				return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
			}
			// create user account
			accountRepo := repository.NewUserAccountRepository()

			userAccout := &models.UserAccountModel{
				UserId:         existUser.ID,
				OrganizationId: &user.Accessor,
				Active:         true,
			}
			account, err := accountRepo.CreateUserAccount(userAccout)
			if err != nil || account == nil {
				return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user account")
			}

			// create user role
			roleRepo := repository.NewUserRoleRepository()
			userRole := &models.UserRoleModel{
				UserId:         existUser.ID,
				OrganizationId: user.Accessor,
				RoleId:         payload.RoleId,
			}
			role, err := roleRepo.CreateUserRole(userRole)
			if err != nil || role == nil {
				return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user role")
			}

			// add staff role to user
			return c.Status(fiber.StatusOK).JSON(existUser)

		}

		// if user is already a staff in the organization

		return utils.WriteError(c, fiber.StatusBadRequest, "email already exist")
	}

	password, _ := utils.GeneratePassword(8)
	country := strings.ToUpper(payload.Country)

	newStaff := &models.UserModel{
		FirstName:  payload.FirstName,
		MiddleName: payload.MiddleName,
		LastName:   payload.LastName,
		Email:      strings.ToLower(payload.Email),
		Gender:     payload.Gender,
		Phone:      &payload.Phone,
		Country:    &country,
		Password:   utils.HashPassword([]byte(password)),
	}

	staff, err := userRepo.CreateUser(newStaff)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
	}

	// create organization user
	newOrgStaff := &models.OrganizationUserModel{
		OrganizationId: user.Accessor,
		UserId:         staff.ID,
		CreatedById:    &user.OwnerID,
	}
	orgUser, err := orgUserRepo.CreateOrganizationUser(newOrgStaff)
	if err != nil || orgUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create organization user")
	}

	// create user account
	accountRepo := repository.NewUserAccountRepository()

	userAccout := &models.UserAccountModel{
		UserId:         staff.ID,
		OrganizationId: &user.Accessor,
		Active:         true,
	}
	account, err := accountRepo.CreateUserAccount(userAccout)
	if err != nil || account == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user account")
	}

	// create user role
	roleRepo := repository.NewUserRoleRepository()
	userRole := &models.UserRoleModel{
		UserId:         staff.ID,
		OrganizationId: user.Accessor,
		RoleId:         payload.RoleId,
	}
	role, err := roleRepo.CreateUserRole(userRole)
	if err != nil || role == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user role")
	}

	organization, err := orgRepo.GetOrganizationById(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
	}

	email := types.EmailPayload{
		Name:         staff.FirstName + " " + staff.LastName,
		MailTo:       staff.Email,
		Subject:      "New Account " + organization.Name,
		Body:         template.HTML("<p>You have been added to " + organization.Name + ".</p><p>Use the below password to access your account and change your password.</p><p>Should you need any assistance or have questions, our support team is always ready to help. You can reach out to us at any time, and we'll ensure you receive the support you need.</p><p>Your password is <b>" + password + "</b>.</p>"),
		Code:         &password,
		TemplateFile: "registration.html",
	}

	// send email to the user
	go emails.SendEmail(email)

	return c.Status(fiber.StatusOK).JSON(staff)

}

func OrganizationStaffs(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	orgUserSummaryRepo := repository.NewOrganizationUserSummaryRepository()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")
	data, totalCount, err := orgUserSummaryRepo.GetOrganizationUserSummaryPageable(user.Accessor, page, pageSize, search)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func GetOrganizationStaffById(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	orgUserRepo := repository.NewOrganizationUserSummaryRepository()
	orgUser, err := orgUserRepo.GetOrganizationUserSummaryByUserId(uuid.MustParse(c.Params("staffId")), user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(orgUser)

}

func OrganizationUserAnalytics(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	orgUserAnalyticsRepo := repository.NewOrganizationUserAnalyticsRepository()
	analytics, err := orgUserAnalyticsRepo.GetOrganizationUserAnalyticsByOrganizationId(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(analytics)

}

func GetOrganizationRoles(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	roleRepo := repository.NewRoleRepository()
	roles, err := roleRepo.GetRolesByOrganizationId(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(roles)

}

func CreateOrganizationRoles(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	request := new(types.CreateRoleRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "bad request")
	}

	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}
	roleRepo := repository.NewRoleRepository()

	// check if role already exist in the organization
	existRole, err := roleRepo.ExistsRoleByNameAndOrganizationId(request.Name, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if existRole {
		return utils.WriteError(c, fiber.StatusBadRequest, "role already exist")
	}

	role := &models.RoleModel{
		Name:           request.Name,
		OrganizationId: &user.Accessor,
	}

	createRole, err := roleRepo.CreateRole(role)

	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	return c.Status(fiber.StatusCreated).JSON(createRole)

}

func UpdateOrganizationRoles(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	request := new(types.UpdateRoleRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "bad request")
	}

	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}
	roleRepo := repository.NewRoleRepository()

	// find the role
	role, err := roleRepo.GetRoleByIdAndOrganizationId(request.RoleId, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if role == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "role not found")
	}

	role.Name = request.Name

	updated, err := roleRepo.UpdateRole(role)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(updated)

}

func GetPermissions(c *fiber.Ctx) error {
	// user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	// if !utils.CheckAccessPermission("Manage Permissions", user.Permissions) {
	// 	return utils.WriteError(c, fiber.StatusForbidden, "Acccess denied")
	// }
	permissionRepo := repository.NewPermissionRepository()
	permissions, err := permissionRepo.GetPermissionsByAccount("ORGANIZATION")
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(permissions)

}

func GetStaffPermissions(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	userAccRepo := repository.NewUserAccountRepository()
	userAccPermRepo := repository.NewUserAccountPermissionRepository()

	account, _accErr := userAccRepo.GetUserAccountByUserIdAndOrganizationId(uuid.MustParse(c.Params("staffId")), user.Accessor)
	if _accErr != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}
	if account == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}

	perms, errPerm := userAccPermRepo.GetPermissionsByUserIdAndAccountId(user.OwnerID, account.ID)

	if errPerm != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to fetch permissions")
	}

	return c.Status(fiber.StatusOK).JSON(perms)

}

func UpdateStaffPermissions(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	request := new(types.UpdateStaffPermissionsRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "bad request")
	}

	staffId := uuid.MustParse(c.Params("staffId"))

	userAccRepo := repository.NewUserAccountRepository()
	userAccPermRepo := repository.NewUserAccountPermissionRepository()

	account, _accErr := userAccRepo.GetUserAccountByUserIdAndOrganizationId(staffId, user.Accessor)
	if _accErr != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}
	if account == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}

	// remove all permissions from the account
	err := userAccPermRepo.DeleteUserAccountPermissionsByUserIdAndAccountId(user.OwnerID, account.ID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to update permissions")
	}

	// loop through the permissions and check if they exist and add them to user account
	for _, permission := range request.PermissionIds {
		permission, err := repository.NewPermissionRepository().GetPermissionById(permission)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to update permissions")
		}
		if permission == nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "Permission not found")
		}
		userAccPermRepo.CreateUserAccountPermission(&models.UserAccountPermissionModel{
			OrganizationId: user.Accessor,
			UserId:         user.OwnerID,
			AccountId:      account.ID,
			PermissionId:   permission.ID,
		})
	}

	userPermissions, err := userAccPermRepo.GetPermissionsByUserIdAndAccountId(user.OwnerID, account.ID)

	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to update permissions")
	}
	return c.Status(fiber.StatusOK).JSON(userPermissions)

}
