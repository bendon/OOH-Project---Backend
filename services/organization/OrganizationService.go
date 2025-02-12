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

	existUser, err := userRepo.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create staff")
	}
	if existUser != nil {
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
	orgUserRepo := repository.NewOrganizationUserRepository()

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
