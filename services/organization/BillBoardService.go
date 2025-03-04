package services

import (
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	emails "bbscout/services/email"
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
		Description:     &payload.Description,
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
		ObjectType:      payload.ObjectType,
		OwnerContacts:   payload.OwnerContacts,
		Owner:           payload.Owner,
		OwnerEmails:     payload.OwnerEmail,
		Occupied:        payload.Occupied,
		City:            payload.City,
		CloseUpImageId:  payload.CloseUpImageId,
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
	billboard.Description = &payload.Description
	billboard.ImageId = &file.ID
	billboard.Location = payload.Location
	billboard.Latitude = payload.Latitude
	billboard.Longitude = payload.Longitude
	billboard.Width = payload.Width
	billboard.Height = payload.Height
	billboard.Unit = payload.Unit
	billboard.Type = payload.Type
	billboard.Price = &payload.Price
	billboard.ObjectType = payload.ObjectType
	billboard.OwnerContacts = payload.OwnerContacts
	billboard.Owner = payload.Owner
	billboard.OwnerEmails = payload.OwnerEmail
	billboard.Occupied = payload.Occupied
	billboard.City = payload.City
	billboard.CloseUpImageId = payload.CloseUpImageId

	updatedBillBoard, err := billboardRepo.UpdateBillBoard(billboard)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusOK).JSON(updatedBillBoard)
}

func CreateBillboardCampaign(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	var payload types.CreateBillboardCampaignRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}
	//validate request
	if err := validate.Struct(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	// check if the billboard id has active campaign
	billboardCampRepo := repository.NewBillboardCampaignRepository()
	active, err := billboardCampRepo.FindBillboardCampaignByOrganizationIdAndBillboardIdAndActive(user.Accessor, payload.BillboardId, true)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if active != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "bill board has active campaign")
	}

	//// check if the billboard id exists
	billboardRepo := repository.NewBillBoardRepository()
	billboard, err := billboardRepo.GetBillBoardByIdAndOrganizationId(payload.BillboardId, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if billboard == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "billboard not found")
	}
	// create the campaign
	var startDate *int64
	var endDate *int64

	// check if the payload has start date and end date thenconvert them to unix
	if payload.StartDate != nil {
		date, err := time.Parse("2006-01-02", *payload.StartDate)
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid start date")
		}
		unix := date.Unix()
		startDate = &unix
	}
	if payload.EndDate != nil {
		date, err := time.Parse("2006-01-02", *payload.EndDate)
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid end date")
		}
		unix := date.Unix()
		endDate = &unix
	}

	campaign := &models.BillboardCampaignModel{
		CampaignBrand:       payload.CampaignBrand,
		Others:              payload.OtherDetails,
		OrganizationId:      user.Accessor,
		BillboardId:         payload.BillboardId,
		StartDate:           startDate,
		EndDate:             endDate,
		CampaignDescription: payload.CampaignDescription,
		CampaignInsights:    payload.CampaignInsight,
		CreatedById:         user.OwnerID,
		Location:            payload.Location,
		ClientFirstName:     payload.ClientFirstName,
		ClientLastName:      payload.ClientLastName,
		Email:               payload.Email,
		Phone:               payload.Phone,
		ImageId:             payload.ImageId,
		Active:              true,
		TargetAudience:      payload.TargetAudience,
		TargetAge:           payload.TargetAge,
		TargetGender:        payload.TargetGender,
		CampaignSocials:     payload.CampaignSocials,
		Products:            payload.Products,
		SiteUrl:             payload.SiteUrl,
	}

	createdCampaign, err := billboardCampRepo.CreateBillboardCampaign(campaign)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusCreated).JSON(createdCampaign)

}

func GetBillBaordsOrganizationCampaigns(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardCampRepo := repository.NewBillboardCampaignRepository()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")
	data, totalCount, err := billboardCampRepo.FindBillboardCampaignByOrganizationId(user.Accessor, page, pageSize, search)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)
}

func GetCampaignById(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardCampRepo := repository.NewBillboardCampaignRepository()
	id := uuid.MustParse(c.Params("campaignId"))
	campaign, err := billboardCampRepo.GetBillboardCampaignByIdAndOrganizationId(id, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if campaign == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "campaign not found")
	}
	return c.Status(fiber.StatusOK).JSON(campaign)

}

func GetBillboardHistryCampaigns(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardCampRepo := repository.NewBillboardCampaignRepository()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")
	data, totalCount, err := billboardCampRepo.FindBillboardCampaignByOrganizationIdAndBillboardIdPageable(user.Accessor, uuid.MustParse(c.Params("billboardId")), page, pageSize, search)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func UpdateBillboardCampaign(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	payload := new(types.UpdateBillboardCampaignRequest)
	if err := c.BodyParser(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	billboardCampRepo := repository.NewBillboardCampaignRepository()

	compaignId := uuid.MustParse(c.Params("campaignId"))

	// validate the payload
	if err := validate.Struct(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	// check if the campaign id exists
	campaign, err := billboardCampRepo.GetBillboardCampaignByIdAndOrganizationId(compaignId, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if campaign == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "campaign not found")
	}

	// update the campaign
	var startDate *int64
	var endDate *int64

	// check if the payload has start date and end date thenconvert them to unix
	if payload.StartDate != nil {
		date, err := time.Parse("2006-01-02", *payload.StartDate)
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid start date")
		}
		unix := date.Unix()
		startDate = &unix
	}
	if payload.EndDate != nil {
		date, err := time.Parse("2006-01-02", *payload.EndDate)
		if err != nil {
			return utils.WriteError(c, fiber.StatusBadRequest, "invalid end date")
		}
		unix := date.Unix()
		endDate = &unix
	}

	campaign.StartDate = startDate
	campaign.EndDate = endDate
	campaign.CampaignDescription = payload.CampaignDescription
	campaign.CampaignInsights = payload.CampaignInsight
	campaign.Location = payload.Location
	campaign.ClientFirstName = payload.ClientFirstName
	campaign.ClientLastName = payload.ClientLastName
	campaign.Email = payload.Email
	campaign.Phone = payload.Phone
	campaign.ImageId = payload.ImageId
	campaign.CampaignBrand = payload.CampaignBrand
	campaign.Others = payload.Others
	campaign.TargetAge = payload.TargetAge
	campaign.TargetGender = payload.TargetGender
	campaign.Products = payload.Products
	campaign.SiteUrl = payload.SiteUrl
	campaign.TargetAudience = payload.TargetAudience
	updated, err := billboardCampRepo.UpdateBillboardCampaign(campaign)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	return c.Status(fiber.StatusOK).JSON(updated)

}
func DeleteBillboardCampaign(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardCampRepo := repository.NewBillboardCampaignRepository()
	compaignId := uuid.MustParse(c.Params("campaignId"))

	// check if the campaign id exists
	campaign, err := billboardCampRepo.GetBillboardCampaignByIdAndOrganizationId(compaignId, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if campaign == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "campaign not found")
	}

	// update the campaign
	campaign.Active = false
	updated, err := billboardCampRepo.UpdateBillboardCampaign(campaign)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	// delete the campaign

	err = billboardCampRepo.DeleteBillboardCampaign(updated.ID)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)

}

func CloseBillboardCampaign(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardCampRepo := repository.NewBillboardCampaignRepository()
	compaignId := uuid.MustParse(c.Params("campaignId"))

	// check if the campaign id exists
	campaign, err := billboardCampRepo.GetBillboardCampaignByIdAndOrganizationId(compaignId, user.Accessor)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if campaign == nil {
		return utils.WriteError(c, fiber.StatusNotFound, "campaign not found")
	}

	// update the campaign close
	closeDate := time.Now().Unix()
	campaign.Active = false
	campaign.ClosedDate = &closeDate
	updated, err := billboardCampRepo.UpdateBillboardCampaign(campaign)

	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	return c.Status(fiber.StatusOK).JSON(updated)

}

func MyBillBoardsUploads(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardSummaryRepo := repository.NewBillBoardSummaryRepository()

	userId := user.OwnerID
	if c.Query("userId", "") != "" {
		userId = uuid.MustParse(c.Query("userId"))

	}
	// get staff uploads
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("search", "")

	data, totalCount, err := billboardSummaryRepo.GetStaffBillBoardsSummary(user.Accessor, userId, page, pageSize, search)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func GetBillboardDailyFilter(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardSummaryRepo := repository.NewBillBoardSummaryRepository()

	// Define time format
	layout := "2006-01-02 15:04:05"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
	search := c.Query("code", "")
	//getstart and end date from query other pick todays date
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02 15:04:05")
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).Format("2006-01-02 15:04:05")

	if c.Query("startDate", "") != "" {
		startDate = c.Query("startDate") + " 00:00:00"

	}
	if c.Query("endDate", "") != "" {
		endDate = c.Query("endDate") + " 23:59:59"
	}

	startTime, _ := time.Parse(layout, startDate)
	endTime, _ := time.Parse(layout, endDate)

	// convert them to unix
	startUnix := startTime.Unix()
	endUnix := endTime.Unix()

	data, totalCount, err := billboardSummaryRepo.GetBillboardDailyFilterPageable(user.Accessor, startUnix, endUnix, page, pageSize, search)

	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}

func MyBillBoardById(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardSummaryRepo := repository.NewBillBoardSummaryRepository()

	id := uuid.MustParse(c.Params("billboardId"))

	board, err := billboardSummaryRepo.GetStaffBillBoardsSummaryById(id, user.OwnerID)

	if err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error fetching the billboard")
	}

	return c.Status(fiber.StatusOK).JSON(board)
}

func SendEmail(c *fiber.Ctx) error {

	request := new(types.SendEmailRequest)
	if err := c.BodyParser(request); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}

	email := types.EmailPayload{
		Name:         request.Name,
		MailTo:       request.Email,
		Subject:      request.Subject,
		Body:         template.HTML(request.Body),
		Link:         request.Link,
		TemplateFile: "registration.html",
	}

	// send email to the user
	go emails.SendEmail(email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "email sent successfully",
	})

}
func CreateBillboardTypes(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	payload := new(types.BillboardTypeRequest)
	if err := c.BodyParser(payload); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid request")
	}

	//check if the billboard type already exists
	billboardTypeRepo := repository.NewBillboardTypesRepository()
	exists, err := billboardTypeRepo.ExistsBillboardTypeByName(strings.ToUpper(payload.Name))
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if exists {
		return utils.WriteError(c, fiber.StatusBadRequest, "billboard type already exists")
	}

	// create billboard type

	types := &models.BillboardTypesModel{
		Name:        strings.ToUpper(payload.Name),
		CreatedById: user.OwnerID,
	}

	created, err := billboardTypeRepo.CreateBillboardType(types)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusCreated).JSON(created)

}

func DeleteBillboardType(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	billboardTypeRepo := repository.NewBillboardTypesRepository()

	id := uuid.MustParse(c.Params("typeId"))

	exists, err := billboardTypeRepo.ExistsBillboardTypeById(id)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	if !exists {
		return utils.WriteError(c, fiber.StatusNotFound, "billboard type not found")
	}

	// updated deleted by

	boardTYpe, err := billboardTypeRepo.GetBillboardTypeByID(id)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	boardTYpe.DeletedById = &user.OwnerID

	billboardTypeRepo.UpdateBillboardType(boardTYpe)

	// delete billboard type
	err = billboardTypeRepo.DeleteBillboardType(id)
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}
	return c.Status(fiber.StatusNoContent).JSON(nil)

}

func GetBillboardTypes(c *fiber.Ctx) error {
	billboardTypeRepo := repository.NewBillboardTypesRepository()

	types, err := billboardTypeRepo.GetBillboardTypes()
	if err != nil {
		return utils.WriteError(c, fiber.StatusInternalServerError, "server error")
	}

	return c.Status(fiber.StatusOK).JSON(types)

}
