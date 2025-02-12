package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	random "math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

func CheckIfUUID(n string) (uuid.UUID, error) {
	id, err := uuid.Parse(n)
	return id, err
}

func WriteError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
	})
}

func HashPassword(password []byte) string {

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(password []byte, hashed string) bool {
	byteHash := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(byteHash, password)

	return err == nil
}

func CheckAccessPermission(permission string, permissions []string) bool {

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func GeneratePassword(length int) (string, error) {
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordChars))))
		if err != nil {
			return "", err
		}
		password[i] = passwordChars[randomIndex.Int64()]
	}
	return string(password), nil
}

func GenerateRandomThreeLetters() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // Alphabet
	result := make([]byte, 3)                    // Create a slice to hold 3 random letters
	for i := range result {
		result[i] = letters[random.IntN(len(letters))] // Pick a random letter
	}
	return string(result)
}

func ReverseIP(letterCode string) string {
	letterMap := map[rune]string{
		'N': "1",
		'O': "2",
		'P': "3",
		'Q': "4",
		'R': "5",
		'S': "6",
		'T': "7",
		'U': "8",
		'V': "9",
		'X': "0",
		'A': ".",
	}

	result := strings.Split(letterCode, "b47d3")
	var octet string

	for _, letter := range result[1] {
		octet += letterMap[letter]
	}

	return octet
}

func GenerateQRCodeBase64(invoiceNumber string) string {
	// Generate QR code and save as file
	qrPath := generateQRCode(invoiceNumber)

	// Convert to base64
	base64QR, err := encodeToBase64(qrPath)
	if err != nil {
		log.Fatal(err)
	}

	return base64QR
}

func encodeToBase64(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func generateQRCode(invoiceNumber string) string {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads", 0755)
		if err != nil {
			fmt.Println("Failed to create uploads directory")
		}
	}
	if _, err := os.Stat("./uploads/qr"); os.IsNotExist(err) {
		err := os.Mkdir("./uploads/qr", 0755)
		if err != nil {
			fmt.Println("Failed to create qr directory")
		}
	}

	name := "uploads/qr/qr_code" + invoiceNumber + ".png"
	// Create the QR code
	err := qrcode.WriteFile(invoiceNumber, qrcode.Medium, 256, name)
	if err != nil {
		log.Fatal(err)
	}

	// Returning file path or base64 encoding
	return name
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

func NewPaginationResponse(data interface{}, total int64, page int, pageSize int) *PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func GetDateTimeByDate(date string) string {
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}
	return parsedTime.Format("2006-01-02")

}
