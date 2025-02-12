package services

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"bbscout/middleware"
	"bbscout/repository"
	"bbscout/types"
	"bbscout/utils"
)

func Login(c *fiber.Ctx) error {
	credentials := new(types.Credentials)
	c.BodyParser(credentials)
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByEmail(credentials.Email)

	if err != nil || user == nil {
		return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized. invalid credentials")
	}

	compared := utils.ComparePasswords([]byte(credentials.Password), user.Password)
	if !compared {
		return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized. invalid credentials")
	}

	expirationTime := jwt.NewNumericDate(time.Now().Add(6 * time.Hour))

	claims := &middleware.Claims{
		OwnerID:  user.ID,
		Username: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)

	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "unauthorized. server error")
	}

	response := middleware.TokenResponse{
		User:        *user,
		AccessToken: tokenString,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}

func GetUserAccounts(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.Claims)
	userAccountRepo := repository.NewUserAccountRepository()
	profile, err := userAccountRepo.GetUserAccountsByUserId(user.OwnerID)

	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "No accounts found for user")
	}

	return c.Status(fiber.StatusOK).JSON(profile)

}

func GetUserProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.Claims)
	userRepo := repository.NewUserRepository()
	profile, err := userRepo.GetUserById(user.OwnerID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "No accounts found for user")
	}
	return c.Status(fiber.StatusOK).JSON(profile)

}

func PostSwtichAccounts(c *fiber.Ctx) error {

	switchAccount := new(types.SwitchAccountPayload)
	c.BodyParser(switchAccount)

	user := c.Locals("user").(*middleware.Claims)
	userRepo := repository.NewUserRepository()
	profile, err := userRepo.GetUserById(user.OwnerID)

	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "No accounts found for user")
	}

	userAccRepo := repository.NewUserAccountRepository()
	userAccPermRepo := repository.NewUserAccountPermissionRepository()

	//check if the account is active or

	account, _accErr := userAccRepo.GetUserAccountByUserIdAndId(profile.ID, switchAccount.AccountId)
	if _accErr != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}

	perms, errPerm := userAccPermRepo.GetPermissionsByUserIdAndAccountId(user.OwnerID, account.ID)

	if errPerm != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to fetch permissions")
	}

	// Extract permission names into a slice of strings
	permissionNames := make([]string, len(perms))
	for i, perm := range perms {
		permissionNames[i] = perm.Name // Assuming PermissionModel has a Name field
	}

	// create new token account
	expirationTime := jwt.NewNumericDate(time.Now().Add(6 * time.Hour))

	claims := &middleware.AccountClaims{
		OwnerID:     profile.ID,
		Username:    profile.Email,
		Accessing:   account.ID,
		Accessor:    account.OrganizationId,
		Permissions: permissionNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)

	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "unauthorized. server error")
	}

	response := middleware.TokenResponse{
		User:         *profile,
		Account:      account,
		AccessToken:  tokenString,
		RefreshToken: tokenString,
		Permissions:  &permissionNames,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
