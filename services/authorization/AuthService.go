package services

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	emails "bbscout/services/email"
	"bbscout/types"
	"bbscout/utils"
)

var validate = validator.New()

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

	if !user.Active {
		user.Active = true
		userRepo.UpdateUser(user)
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
	if account == nil {
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

func RefreshToken(c *fiber.Ctx) error {

	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	userRepo := repository.NewUserRepository()
	profile, err := userRepo.GetUserById(user.OwnerID)

	if err != nil {
		return utils.WriteError(c, fiber.StatusNotFound, "No accounts found for user")
	}

	userAccRepo := repository.NewUserAccountRepository()
	userAccPermRepo := repository.NewUserAccountPermissionRepository()

	//check if the account is active or

	account, _accErr := userAccRepo.GetUserAccountByUserIdAndId(profile.ID, user.Accessing)
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
	expirationTime := jwt.NewNumericDate(time.Now().Add(48 * time.Hour))
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
	userAccRepo := repository.NewUserAccountRepository()
	// // Convert payload to JSON
	userInfo := &types.UserInfo{
		Sub:           payload.Claims["sub"].(string),
		Name:          payload.Claims["name"].(string),
		Email:         payload.Claims["email"].(string),
		EmailVerified: payload.Claims["email_verified"].(bool),
		Picture:       payload.Claims["picture"].(string),
		GivenName:     payload.Claims["given_name"].(string),
		FamilyName:    payload.Claims["family_name"].(string),
	}

	// generate a JWT token

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByEmail(userInfo.Email)

	if err != nil || user == nil {
		password, _ := utils.GeneratePassword(8)
		// create a new user
		newUser := &models.UserModel{
			FirstName:  userInfo.GivenName,
			LastName:   userInfo.FamilyName,
			Email:      userInfo.Email,
			MiddleName: "",
			Password:   utils.HashPassword([]byte(password)),
			Phone:      nil,
			Gender:     1,
			Verified:   userInfo.EmailVerified,
			Active:     true,
		}
		// save the user to the database
		createdUser, err := userRepo.CreateUser(newUser)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Error creating user")
		}

		// create a new userserAccRepo := repository.NewUserAccountRepository()
		roleRepo := repository.NewRoleRepository()
		userAccPermRepo := repository.NewUserAccountPermissionRepository()
		orgUserRepo := repository.NewOrganizationUserRepository()
		organizationRepo := repository.NewOrganizationRepository()
		organization, err := organizationRepo.FindOrganizationLatest()
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting organization")
		}
		if organization == nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting organization")
		}

		role, err := roleRepo.GetRoleByNameAndOrganizationId("USER", organization.ID)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting role")
		}
		if role == nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting role")
		}

		// create organization user
		newOrgStaff := &models.OrganizationUserModel{
			OrganizationId: organization.ID,
			UserId:         createdUser.ID,
			CreatedById:    nil,
		}
		orgUser, err := orgUserRepo.CreateOrganizationUser(newOrgStaff)
		if err != nil || orgUser == nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create organization user")
		}

		// create user account
		accountRepo := repository.NewUserAccountRepository()

		userAccout := &models.UserAccountModel{
			UserId:         createdUser.ID,
			OrganizationId: &organization.ID,
			Active:         true,
		}
		createdAccount, err := accountRepo.CreateUserAccount(userAccout)
		if err != nil || createdAccount == nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user account")
		}

		account, _accErr := userAccRepo.GetUserAccountByUserIdAndId(createdUser.ID, createdAccount.ID)
		if _accErr != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
		}
		if account == nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
		}

		// create user role
		roleUserRepo := repository.NewUserRoleRepository()
		userRole := &models.UserRoleModel{
			UserId:         createdUser.ID,
			OrganizationId: organization.ID,
			RoleId:         role.ID,
		}
		createdRoleUser, err := roleUserRepo.CreateUserRole(userRole)
		if err != nil || createdRoleUser == nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user role")
		}

		perms, errPerm := userAccPermRepo.GetPermissionsByUserIdAndAccountId(createdUser.ID, account.ID)

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
			OwnerID:     createdUser.ID,
			Username:    createdUser.Email,
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
			User:         *createdUser,
			Account:      account,
			AccessToken:  tokenString,
			RefreshToken: *GetRefreshToken(*claims),
			Permissions:  &permissionNames,
		}

		email := types.EmailPayload{
			Name:         createdUser.FirstName + " " + createdUser.LastName,
			MailTo:       createdUser.Email,
			Subject:      "Account Creation On " + organization.Name,
			Body:         template.HTML("<p>Thank you for joining " + organization.Name + ".</p><p>Your account has been created successfully.</p><p>Should you need any assistance or have questions, our support team is always ready to help. You can reach out to us at any time, and we'll ensure you receive the support you need.</p>"),
			TemplateFile: "registration.html",
		}

		// send email to the user
		go emails.SendEmail(email)

		return c.Status(fiber.StatusOK).JSON(response)

	}

	if !user.Active {
		user.Active = true
		userRepo.UpdateUser(user)
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

func Register(c *fiber.Ctx) error {
	// parse the request body
	request := new(types.UserRegisterRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}
	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}
	userAccRepo := repository.NewUserAccountRepository()
	roleRepo := repository.NewRoleRepository()
	userAccPermRepo := repository.NewUserAccountPermissionRepository()
	orgUserRepo := repository.NewOrganizationUserRepository()
	organizationRepo := repository.NewOrganizationRepository()
	organization, err := organizationRepo.FindOrganizationLatest()
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting organization")
	}
	if organization == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting organization")
	}

	role, err := roleRepo.GetRoleByNameAndOrganizationId("USER", organization.ID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting role")
	}
	if role == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting role")
	}

	// check if user already exists
	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetUserByEmail(request.Email)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error getting user")
	}
	if user != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User already exists")
	}
	var phoneNumber *int
	var gender = 1
	if request.Phone != nil {
		phoneNumber = request.Phone

	}

	if request.Gender != nil {
		gender = *request.Gender

	}
	// create a new user
	newUser := &models.UserModel{
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		MiddleName: request.MiddleName,
		Password:   utils.HashPassword([]byte(request.Password)),
		Phone:      phoneNumber,
		Gender:     gender,
	}
	// save the user to the database
	createdUser, err := userRepo.CreateUser(newUser)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error creating user")
	}

	// create organization user
	newOrgStaff := &models.OrganizationUserModel{
		OrganizationId: organization.ID,
		UserId:         createdUser.ID,
		CreatedById:    nil,
	}
	orgUser, err := orgUserRepo.CreateOrganizationUser(newOrgStaff)
	if err != nil || orgUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create organization user")
	}

	// create user account
	accountRepo := repository.NewUserAccountRepository()

	userAccout := &models.UserAccountModel{
		UserId:         createdUser.ID,
		OrganizationId: &organization.ID,
		Active:         true,
	}
	createdAccount, err := accountRepo.CreateUserAccount(userAccout)
	if err != nil || createdAccount == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user account")
	}

	account, _accErr := userAccRepo.GetUserAccountByUserIdAndId(createdUser.ID, createdAccount.ID)
	if _accErr != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}
	if account == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User Account has either been locked or deactivated")
	}

	// create user role
	roleUserRepo := repository.NewUserRoleRepository()
	userRole := &models.UserRoleModel{
		UserId:         createdUser.ID,
		OrganizationId: organization.ID,
		RoleId:         role.ID,
	}
	createdRoleUser, err := roleUserRepo.CreateUserRole(userRole)
	if err != nil || createdRoleUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create user role")
	}

	perms, errPerm := userAccPermRepo.GetPermissionsByUserIdAndAccountId(createdUser.ID, account.ID)

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
		OwnerID:     createdUser.ID,
		Username:    createdUser.Email,
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
		User:         *createdUser,
		Account:      account,
		AccessToken:  tokenString,
		RefreshToken: *GetRefreshToken(*claims),
		Permissions:  &permissionNames,
	}

	email := types.EmailPayload{
		Name:         createdUser.FirstName + " " + createdUser.LastName,
		MailTo:       createdUser.Email,
		Subject:      "Account Creation On " + organization.Name,
		Body:         template.HTML("<p>Thank you for joining " + organization.Name + ".</p><p>Your account has been created successfully.</p><p>Should you need any assistance or have questions, our support team is always ready to help. You can reach out to us at any time, and we'll ensure you receive the support you need.</p>"),
		TemplateFile: "registration.html",
	}

	// send email to the user
	go emails.SendEmail(email)

	return c.Status(fiber.StatusOK).JSON(response)
}

func ForgotPassword(c *fiber.Ctx) error {
	request := new(types.ForgotPasswordRequest)

	// parse body request
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	userRepo := repository.NewUserRepository()

	user, err := userRepo.GetUserByEmail(request.Email)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User not found")
	}

	if user == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "User with the email does not exist")
	}

	newPassword, ErrPassword := utils.GeneratePassword(8)
	if ErrPassword != nil || newPassword == "" {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid credentials provided")
	}

	change := true
	user.IsChange = &change
	user.Password = utils.HashPassword([]byte(newPassword))

	upatedUser, err := userRepo.UpdateUser(user)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error updating user")
	}
	if upatedUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error resetting user password")
	}

	payloadEmail := types.EmailPayload{
		Name:         user.FirstName + " " + user.LastName,
		Body:         template.HTML("<p>Your account password has been reset successfully.</p><p>Use the provided password to access your account and change your password.</p><p>Should you need any assistance or have questions, our support team is always ready to help. You can reach out to us at any time, and we'll ensure you receive the support you need.</p>"),
		Subject:      "Account Password Reset",
		MailTo:       user.Email,
		Code:         &newPassword,
		TemplateFile: "registration.html",
	}

	go emails.SendEmail(payloadEmail)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "successfully reset user account password",
	})

}

func UpdateUserPassword(c *fiber.Ctx) error {
	user := c.Locals("user").(*middleware.Claims)
	request := new(types.UpdatePasswordPayload)

	// parse body request
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	userRepo := repository.NewUserRepository()
	profile, err := userRepo.GetUserById(user.OwnerID)
	if err != nil || profile == nil {
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized access")
	}

	change := false
	profile.IsChange = &change
	profile.Password = utils.HashPassword([]byte(request.Password))
	updatedUser, err := userRepo.UpdateUser(profile)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error updating user")
	}
	if updatedUser == nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Error updating user password")
	}

	payloadEmail := types.EmailPayload{
		Name:         profile.FirstName + " " + profile.LastName,
		Body:         template.HTML("<p>Your account password has been updated successfully.</p><p>Should you need any assistance or have questions, our support team is always ready to help. You can reach out to us at any time, and we'll ensure you receive the support you need.</p>"),
		Subject:      "Account Password Update",
		MailTo:       profile.Email,
		TemplateFile: "registration.html",
	}

	go emails.SendEmail(payloadEmail)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "password updated successfully",
	})
}
