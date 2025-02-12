package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"

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

	tokenCodesl := utils.HashPassword([]byte("token"))

	claims := &middleware.AccountClaims{
		OwnerID:     profile.ID,
		Username:    profile.Email,
		Accessing:   account.ID,
		Accessor:    account.OrganizationId,
		CodeSl:      tokenCodesl,
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
		RefreshToken: *GetRefreshToken(*claims),
		Permissions:  &permissionNames,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetRefreshToken(claims middleware.AccountClaims) *string {
	expirationTime := jwt.NewNumericDate(time.Now().Add(8 * time.Hour))
	refreshTokenCodesl := utils.HashPassword([]byte("refreshToken"))

	claims.CodeSl = refreshTokenCodesl
	claims.RegisteredClaims.ExpiresAt = expirationTime

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)

	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return nil
	}

	return &tokenString
}

func ChangePassword(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	request := new(types.ChangePasswordRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}
	userRepo := repository.NewUserRepository()
	profile, err := userRepo.GetUserById(user.OwnerID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "No accounts found for user")
	}
	// if user not found
	if profile == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "User not found")
	}

	// check if old password is correct
	compared := utils.ComparePasswords([]byte(request.OldPassword), profile.Password)
	if !compared {
		return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized. invalid credentials")
	}

	profile.Password = utils.HashPassword([]byte(request.NewPassword))
	updatedUser, err := userRepo.UpdateUser(profile)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error updating user")
	}
	if updatedUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error updating user")
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
		"user":    updatedUser,
	})

}

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Server is up and running",
	})

}

func AuthGoogleVerify(c *fiber.Ctx) error {
	// parse the request body
	request := new(types.AuthGoogleVerificationRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// Verify token using Google's API
	payload, err := idtoken.Validate(context.Background(), request.Token, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// fmt.Println("Payload:", payload)

	// // Convert payload to JSON
	// userInfo := types.UserInfo{
	// 	Sub:           payload.Claims["sub"].(string),
	// 	Name:          payload.Claims["name"].(string),
	// 	Email:         payload.Claims["email"].(string),
	// 	EmailVerified: payload.Claims["email_verified"].(bool),
	// 	Picture:       payload.Claims["picture"].(string),
	// 	GivenName:     payload.Claims["given_name"].(string),
	// 	FamilyName:    payload.Claims["family_name"].(string),
	// 	Locale:        payload.Claims["locale"].(string),
	// }

	// generate a JWT token

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": payload.Claims,
	})

}
