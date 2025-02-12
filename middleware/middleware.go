package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"bbscout/repository"
	"bbscout/utils"
)

func CheckAuthentication(c *fiber.Ctx) error {
	// Retrieve the token from the "token" cookie
	authHeader := c.Get("Authorization")

	// Check if the token is provided
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized Access",
		})
	}

	// The token should be in the format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")

	}

	// Get the token value from the cookie
	tokenStr := tokenParts[1]
	claims := &Claims{}

	// Parse and validate the token
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")
		}
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")
	}

	if !tkn.Valid {
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	c.Locals("user", claims)

	return c.Next()

}

func CheckAccountAuthentication(c *fiber.Ctx) error {
	// Retrieve the token from the "token" cookie
	authHeader := c.Get("Authorization")

	// Check if the token is provided
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized Access",
		})
	}

	// The token should be in the format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")

	}

	// Get the token value from the cookie
	tokenStr := tokenParts[1]
	claims := &AccountClaims{}

	// Parse and validate the token
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")
		}
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")
	}

	if !tkn.Valid {
		return utils.WriteError(c, fiber.StatusUnauthorized, "Unauthorized Access")
	}

	userAccRepo := repository.NewUserAccountRepository()
	//check active account
	acc, _accErr := userAccRepo.GetUserAccountByUserIdAndId(claims.OwnerID, claims.Accessing)
	if _accErr != nil || acc == nil {
		return utils.WriteError(c, fiber.StatusForbidden, "Access denied")
	}

	if acc.OrganizationId == nil {
		return utils.WriteError(c, fiber.StatusForbidden, "Access denied")
	}

	accPermRepo := repository.NewUserAccountPermissionRepository()

	perms, errPerm := accPermRepo.GetPermissionsByUserIdAndAccountId(claims.OwnerID, claims.Accessing)

	if errPerm != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Unable to fetch permissions")
	}

	// check if the token is not a refresh token
	compared := utils.ComparePasswords([]byte("token"), claims.CodeSl)
	if !compared {
		return utils.WriteError(c, fiber.StatusForbidden, "Access denied")
	}

	// Extract permission names into a slice of strings
	permissionNames := make([]string, len(perms))
	for i, perm := range perms {
		permissionNames[i] = perm.Name // Assuming PermissionModel has a Name field
	}

	// CreateSystemTrail(claims.Username, c)
	// set roles here

	c.Locals("user", c.Locals("user", AccountBranchClaimResponse{
		OwnerID:     claims.OwnerID,
		Accessing:   claims.Accessing,
		Accessor:    *claims.Accessor,
		Username:    claims.Username,
		Permissions: permissionNames,
	}))

	return c.Next()

}

func NotFoundMiddleware(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":   "Cannot " + c.Method() + " " + c.Path(),
		"message": "Route not found",
	})
}
