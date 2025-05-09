package services

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"bbscout/middleware"
	"bbscout/repository"
	"bbscout/types"
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
	billboardMonthlyReportRepo := repository.NewBillboardUploadMonthlyReportRepository()

	now := time.Now()
	currentMonth := now.Month()
	year, _ := now.ISOWeek()
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	monthn, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(currentMonth))))

	report, err := billboardMonthlyReportRepo.GetBillboardUploadMonthlyReport(user.Accessor, yearn, monthn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}
	if report == nil {
		return utils.WriteError(c, fiber.StatusNotFound, " billboard monthly report not found")
	}
	return c.Status(fiber.StatusOK).JSON(report)

}

func BillboardYearlyReports(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	billboardMonthlyReportRepo := repository.NewBillboardUploadMonthlyReportRepository()

	now := time.Now()
	year, _ := now.ISOWeek()
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	reports, err := billboardMonthlyReportRepo.GetBillboardUploadMonthlyReportByYear(user.Accessor, yearn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	return c.Status(fiber.StatusOK).JSON(reports)

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
	var userId uuid.UUID

	if c.Query("userId") != "" {
		id, err := uuid.Parse(c.Query("userId"))
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid user id")
		}
		userId = id
	} else {
		userId = user.OwnerID

	}

	report, err := userUploadWeeklyRepo.GetUserBillboardUploadsReportWeekByUser(user.Accessor, userId, yearn, month, weekn)
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

	var userId uuid.UUID

	if c.Query("userId") != "" {
		id, err := uuid.Parse(c.Query("userId"))
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid user id")
		}
		userId = id
	} else {
		userId = user.OwnerID

	}

	reports, err := userUploadMonthlyRepo.GetUserBillboardUploadMonthReportByUserMonthyly(user.Accessor, userId, yearn, month)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	if reports == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "no report found")

	}

	return c.Status(fiber.StatusOK).JSON(reports)

}

func GetUserUploadReportsYearly(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	userUploadMonthlyRepo := repository.NewUserBillboardUploadMonthReportRepository()

	now := time.Now()
	year, _ := now.ISOWeek()
	yearn, _ := strconv.Atoi(c.Query("year", strconv.Itoa(year)))

	var userId uuid.UUID

	if c.Query("userId") != "" {
		id, err := uuid.Parse(c.Query("userId"))
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid user id")
		}
		userId = id
	} else {
		userId = user.OwnerID

	}

	reports, err := userUploadMonthlyRepo.GetUserBillboardUploadMonthReportByUserYearly(user.Accessor, userId, yearn)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	return c.Status(fiber.StatusOK).JSON(reports)

}

func GetUserGeneralStatistics(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))

	userMonthlyMOnitorRepo := repository.NewUserMonthlyMonitorStatRepository()
	userMonitorRepo := repository.NewUserMonitoringStatRepository()

	userMonthlyAuditRepo := repository.NewUserMonthlyAuditReportRepository()
	userAuditRepo := repository.NewUserAuditReportRepository()

	userMonthlyMonitor, err := userMonthlyMOnitorRepo.GetUserMonthlyMonitorStats(user.Accessor, user.OwnerID, year)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	userMonitor, err := userMonitorRepo.GetUserMonitoringStats(user.Accessor, user.OwnerID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}
	userMonthlyAudit, err := userMonthlyAuditRepo.GetUserMonthlyAuditReport(user.Accessor, user.OwnerID, year)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}
	userAudit, err := userAuditRepo.GetUserAuditReport(user.Accessor, user.OwnerID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "error extracting billboard monthly report")
	}

	monitors := types.MonitorsReportResponse{
		MonthlyReport: userMonthlyMonitor,
		Monitor:       userMonitor,
	}

	audit := types.AuditReportResponse{
		MonthlyReport: userMonthlyAudit,
		Audit:         userAudit,
	}

	report := types.UserReportResponse{
		Monitors: monitors,
		Audit:    audit,
		Year:     year,
	}

	return c.Status(fiber.StatusOK).JSON(report)

}
