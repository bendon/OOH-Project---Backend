package services

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
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

	newStaff := &models.UserModel{
		FirstName:  payload.FirstName,
		MiddleName: payload.MiddleName,
		LastName:   payload.LastName,
		Email:      strings.ToLower(payload.Email),
		Gender:     payload.Gender,
		Phone:      payload.Phone,
		Country:    strings.ToUpper(payload.Country),
		Password:   utils.HashPassword([]byte("password1234")),
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

	return c.Status(fiber.StatusOK).JSON(staff)

}

func OrganizationStaffs(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	orgUserSummaryRepo := repository.NewOrganizationUserSummaryRepository()
	staffs, err := orgUserSummaryRepo.GetOrganizationUserSummary(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(staffs)

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
