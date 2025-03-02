package services

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	"bbscout/repository"
	"bbscout/utils"
)

func BillboardTypeReports(c *fiber.Ctx) error {

	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("type", "")

	billboardTypeReportRepo := repository.NewBillboardOrganizationTypeReportRepository()

	data, totalCount, err := billboardTypeReportRepo.GetBillboardOrganizationTypeReport(user.Accessor, page, pageSize, search)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard type report")
	}

	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)
}

func BillboardOrganizationReports(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	billboardOrganizationReportRepo := repository.NewBillboardOrganizationReportRepository()

	report, err := billboardOrganizationReportRepo.GetBillboardOrganizationReport(user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard organization report")
	}
	return c.Status(fiber.StatusOK).JSON(report)

}

func BillboardLocationReports(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	billboardlocationRepo := repository.NewBillboardOrganizationLocationReportRepository()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("location", "")
	data, totalCount, err := billboardlocationRepo.GetBillboardOrganizationLocationReport(user.Accessor, page, pageSize, search)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard location report")
	}

	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func BillboardWeeklyReports(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	// get currect week number, year and month
	now := time.Now()
	currentMonth := now.Month()
	year, week := now.ISOWeek()
	weekn, _ := strconv.Atoi(c.Query("week", strconv.Itoa(week)))
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(currentMonth))))
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	billboardWeeklyReportRepo := repository.NewBillboardUploadDayOfWeekRepository()

	report, err := billboardWeeklyReportRepo.GetBillboardUploadDayOfWeek(user.Accessor, yearn, month, weekn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard weekly report")
	}
	return c.Status(fiber.StatusOK).JSON(report)

}

func BillboardMonthlyReports(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	month := c.Query("month", "")
	billboardMonthlyReportRepo := repository.NewBillboardUploadMonthlyReportRepository()

	now := time.Now()
	currentMonth := now.Month()
	year, _ := now.ISOWeek()
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	if month == "" {

		reports, err := billboardMonthlyReportRepo.GetBillboardUploadMonthlyReportByYear(user.Accessor, yearn)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
		}

		return c.Status(fiber.StatusOK).JSON(reports)

	}

	monthn, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(currentMonth))))

	report, err := billboardMonthlyReportRepo.GetBillboardUploadMonthlyReport(user.Accessor, yearn, monthn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}
	return c.Status(fiber.StatusOK).JSON(report)

}

func GetMyUploadsSummary(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	userBillboardUploadRepo := repository.NewUserBillboardUploadReportRepository()

	report, err := userBillboardUploadRepo.GetUserBillboardUploadReportByUser(user.Accessor, user.OwnerID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard upload summary")
	}
	return c.Status(fiber.StatusOK).JSON(report)

}

func GetUserOrganizationUploadsSummary(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	userBillboardUploadRepo := repository.NewUserBillboardUploadReportRepository()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("user", "")
	data, totalCount, err := userBillboardUploadRepo.GetUserBillboardUploadReportByOrganizationPageable(user.Accessor, page, pageSize, search)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard location report")
	}

	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func GetUserUploadReportsWeekly(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	userUploadWeeklyRepo := repository.NewUserBillboardUploadsReportWeekRepository()

	now := time.Now()
	currentMonth := now.Month()
	year, week := now.ISOWeek()
	weekn, _ := strconv.Atoi(c.Query("week", strconv.Itoa(week)))
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(currentMonth))))
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	report, err := userUploadWeeklyRepo.GetUserBillboardUploadsReportWeekByUser(user.Accessor, user.OwnerID, yearn, month, weekn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard weekly report")
	}

	return c.Status(fiber.StatusOK).JSON(report)

}

func GetUserUploadReportsMonthly(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	userUploadMonthlyRepo := repository.NewUserBillboardUploadMonthReportRepository()

	now := time.Now()
	currentMonth := now.Month()
	year, _ := now.ISOWeek()
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(currentMonth))))
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	reports, err := userUploadMonthlyRepo.GetUserBillboardUploadMonthReportByUser(user.Accessor, user.OwnerID, yearn, month)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	return c.Status(fiber.StatusOK).JSON(reports)

}
