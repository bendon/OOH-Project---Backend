package services

import (
	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	"bbscout/types"
	"bbscout/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateMonitoringRecord(c *fiber.Ctx) error {

	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	request := new(types.CreateMonitoringRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "Invalid request")
	}

	billboardRepo := repository.NewBillBoardRepository()
	fileRepo := repository.NewFileRepository()

	// validate the request
	if err := validate.Struct(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	// check if the billboard id exists
	if request.BillboardId != nil {
		billboard, err := billboardRepo.GetBillBoardByIdAndOrganizationId(*request.BillboardId, user.Accessor)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
		}
		if billboard == nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "billboard not found")
		}
	}

	// check if long short image id exists
	if request.LongShotImageId != nil {

		// check if the image Id exists
		file, err := fileRepo.GetFileById(*request.LongShotImageId)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
		}
		if file == nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid long short image id")
		}

	}

	// check if close up image id exists
	if request.CloseUpImageId != nil {

		// check if the image Id exists
		file, err := fileRepo.GetFileById(*request.CloseUpImageId)
		if err != nil {
			return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
		}
		if file == nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid close up image id")
		}
	}

	// get today date y-m-d
	today := time.Now()
	date := today.Format("2006-01-02")

	// create monitoring record
	monitoringRepo := repository.NewMonitoringRepository()

	monitor := &models.MonitoringModel{
		BillboardId:          request.BillboardId,
		OrganizationId:       user.Accessor,
		MonitoredById:        user.OwnerID,
		Date:                 &date,
		County:               request.County,
		Street:               request.Street,
		Location:             request.Location,
		Building:             request.Building,
		LongShotImageId:      request.LongShotImageId,
		CloseUpImageId:       request.CloseUpImageId,
		OwnerContacts:        request.OwnerContacts,
		OwnerEmails:          request.OwnerEmail,
		Brand:                request.Brand,
		Campain:              request.Campain,
		Width:                request.Width,
		Height:               request.Height,
		Unit:                 request.Unit,
		Angle:                request.Angle,
		Environment:          request.Environment,
		Illumination:         request.Illumination,
		Material:             request.Material,
		ConditionOfMaterial:  request.ConditionOfMaterial,
		ConditionOfStructure: request.ConditionOfStructure,
		Comments:             request.Comments,
	}

	// create monitoring record
	monitoring, err := monitoringRepo.CreateMonitoring(monitor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to create monitoring record")
	}

	return c.Status(fiber.StatusCreated).JSON(monitoring)

}

func GetMonitoringRecords(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "20"))

	monitoringRepo := repository.NewMonitoringRepository()
	monitoring, total, err := monitoringRepo.GetAllMonitoring(user.Accessor, page, size)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to get monitoring records")
	}

	response := utils.NewPaginationResponse(monitoring, total, page, size)
	return c.Status(fiber.StatusOK).JSON(response)
}
func GetMyMonitoringRecordByUser(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "20"))

	monitoringRepo := repository.NewMonitoringRepository()
	monitoring, total, err := monitoringRepo.GetMonitoringByUser(user.Accessor, user.OwnerID, page, size)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "Failed to get monitoring records")
	}

	response := utils.NewPaginationResponse(monitoring, total, page, size)
	return c.Status(fiber.StatusOK).JSON(response)

}
