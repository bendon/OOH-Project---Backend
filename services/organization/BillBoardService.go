package services

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	"bbscout/types"
	"bbscout/utils"
)

var validate = validator.New()

func CreateBillBoard(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	// parse request body
	var payload types.CreateBillboardRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}

	//validate request
	if err := validate.Struct(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}
	billboardRepo := repository.NewBillBoardRepository()
	fileRepo := repository.NewFileRepository()

	// check if the image Id exists
	file, err := fileRepo.GetFileById(payload.ImageID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if file == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid image id")
	}

	// create billboard sequence
	boardSequenceRepo := repository.NewBillBoardSequenceRepository()

	seqence := &models.BillboardSequenceModel{
		OrganizationId: user.Accessor,
	}
	createdSequence, err := boardSequenceRepo.CreateBillBoardSequence(seqence)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create billboard sequence")
	}
	billboardSequence, err := boardSequenceRepo.GetBillBoardByIdAndOrganizationId(createdSequence.ID, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	boardCode := "BB" + utils.FormatToSixDigits(int(billboardSequence.BoardNumber))

	// Create Billboard instance
	billboard := &models.BillboardModel{
		Description:     payload.Description,
		BoardCode:       boardCode,
		ImageId:         &payload.ImageID,
		Location:        payload.Location,
		Latitude:        payload.Latitude,
		Longitude:       payload.Longitude,
		Accuracy:        payload.Accuracy,
		ParentBoardCode: payload.ParentBoardCode,
		Width:           payload.Width,
		Height:          payload.Height,
		Unit:            payload.Unit,
		Type:            payload.Type,
		Price:           &payload.Price,
		CreatedById:     user.OwnerID,
		OrganizationId:  user.Accessor,
	}

	// create billboard
	newBillBoard, err := billboardRepo.CreateBillBoard(billboard)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "failed to create billboard")
	}
	return c.Status(fiber.StatusCreated).JSON(newBillBoard)

}

func GetBillBoardById(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardRepo := repository.NewBillBoardRepository()
	billboard, err := billboardRepo.GetBillBoardByIdAndOrganizationId(uuid.MustParse(c.Params("id")), user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if billboard == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "billboard not found")
	}
	return c.Status(fiber.StatusOK).JSON(billboard)

}

func GetBillBoards(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardRepo := repository.NewBillBoardRepository()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")

	data, totalCount, err := billboardRepo.GetBillBoardsByOrganizationIdPageable(user.Accessor, page, pageSize, search)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteBillBoard(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardRepo := repository.NewBillBoardRepository()
	billboard, err := billboardRepo.GetBillBoardByIdAndOrganizationId(uuid.MustParse(c.Params("id")), user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if billboard == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "billboard not found")
	}
	err = billboardRepo.DeleteBillBoard(billboard.ID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return utils.WriteError(c, fiber.StatusOK, "billboard deleted successfully")

}

func UpdateBillBoard(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	var payload types.UpdateBillboardRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}
	//validate request
	if err := validate.Struct(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	fileRepo := repository.NewFileRepository()

	// get billboard
	billboardRepo := repository.NewBillBoardRepository()
	billboardOrg, err := billboardRepo.GetBillBoardByIdAndOrganizationId(uuid.MustParse(c.Params("id")), user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if billboardOrg == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "billboard not found")
	}

	// check if the image Id exists
	file, err := fileRepo.GetFileById(payload.ImageID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if file == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid image id")
	}

	billboard, _ := billboardRepo.GetBillBoardById(billboardOrg.ID)

	// update billboard
	billboard.Accuracy = payload.Accuracy
	billboard.ParentBoardCode = payload.ParentBoardCode
	billboard.Description = payload.Description
	billboard.ImageId = &file.ID
	billboard.Location = payload.Location
	billboard.Latitude = payload.Latitude
	billboard.Longitude = payload.Longitude
	billboard.Width = payload.Width
	billboard.Height = payload.Height
	billboard.Unit = payload.Unit
	billboard.Type = payload.Type
	billboard.Price = &payload.Price
	updatedBillBoard, err := billboardRepo.UpdateBillBoard(billboard)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(updatedBillBoard)
}
