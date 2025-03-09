package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"bbscout/middleware"
	"bbscout/models"
	"bbscout/repository"
	"bbscout/services/s3"
	"bbscout/utils"
)

func UploadFile(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)
	fileRepo := repository.NewFileRepository()

	file, err := c.FormFile("file")
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to get file")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to open file")
	}
	defer src.Close()

	fileName := uuid.New().String() + "" + filepath.Ext(file.Filename)

	// Upload to S3 directly from memory
	err = s3.UploadFileToBucket(fileName, src, file)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to upload file to S3")
	}

	// save file details to table
	fileEntry := &models.FileModel{
		FileName:       fileName,
		FileSize:       file.Size,
		FileType:       file.Header.Get("Content-Type"),
		FileExtension:  filepath.Ext(file.Filename),
		FileUrl:        fileName,
		UploadedById:   &user.OwnerID,
		OrganizationId: &user.Accessor,
	}

	created, err := fileRepo.CreateFile(fileEntry)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to save file details")
	}

	return c.Status(200).JSON(created)

}

func GetFileByName(c *fiber.Ctx) error {

	key := c.Params("fileName")

	// get the file from s3
	filePath, err := downloadAndSaveFile(key)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to download file")
	}
	// Serve the file for download
	return c.SendFile(filePath)

}

func downloadAndSaveFile(key string) (string, error) {
	body, err := s3.DownloadFileFromBucket(key)
	if err != nil {
		return "", err
	}
	defer body.Close()

	if _, err := os.Stat("./storage-documents"); os.IsNotExist(err) {
		err := os.Mkdir("./storage-documents", 0755)
		if err != nil {

			fmt.Printf("Error creating directory: %v\n", err)
		}
	}

	// Create a temporary file
	tempDir := "./storage-documents/" // Uses the system temp directory
	tempFilePath := filepath.Join(tempDir, key)

	// Create the file
	outFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Copy content from S3 to the file
	_, err = io.Copy(outFile, body)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil
}

func GetFiles(c *fiber.Ctx) error {
	fileRepo := repository.NewFilesSummaryRepository()
	user := c.Locals("user").(middleware.AccountBranchClaimResponse)

	layout := "2006-01-02 15:04:05"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("size", "20"))
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

	data, totalCount, err := fileRepo.GetFilesSummaryPageable(user.Accessor, startUnix, endUnix, page, pageSize)
	if err != nil || data == nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "error extracting user list")
	}
	response := utils.NewPaginationResponse(data, totalCount, page, pageSize)
	return c.Status(fiber.StatusOK).JSON(response)

}
