package services

import (
	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	"bbscout/repository"
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
