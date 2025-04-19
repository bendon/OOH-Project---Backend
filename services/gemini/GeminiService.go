package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"bbscout/utils"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text       string `json:"text,omitempty"`
	InlineData *Image `json:"inline_data,omitempty"`
}

type Image struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

type Candidate struct {
	Content Content `json:"content"`
}

// Define the structure of the Gemini API response
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Define the expected JSON structure inside the response
type CampaignDetails struct {
	CampaignBrand       string `json:"campaign_brand"`
	CampaignDescription string `json:"campaign_description"`
	ContactInformation  struct {
		Phone []*int    `json:"campaign_phone"`
		Email []*string `json:"campaign_email"`
	} `json:"campaign_contacts"`
	Location              *string `json:"location"`
	BillboardMeasurements struct {
		Height float64 `json:"height"`
		Width  float64 `json:"width"`
		Units  string  `json:"units"`
	} `json:"billboard_measurements"`
	TargetAudience      string      `json:"target_audience"`
	AdditionalNotes     string      `json:"additional_notes"`
	PercentageAccuracy  float64     `json:"percentage_accuracy"`
	SiteUrl             []string    `json:"campaign_site_url"`
	TargetAge           string      `json:"target_age"`
	TargetGender        string      `json:"target_gender"`
	CampaignSocials     interface{} `json:"campaign_socials"`
	OtherDetailts       interface{} `json:"other_details"`
	DetectionConfidence float64     `json:"detection_confidence"`
	ObjectType          *string     `json:"object_type"`
	BillboardType       *string     `json:"billboard_type"`
	Owner               interface{} `json:"owner"`
	Structure           *string     `json:"structure"`
	Material            *string     `json:"material"`
	Illumination        *string     `json:"illumination"`
	Visibility          *string     `json:"visibility"`
	Angle               *string     `json:"angle"`
	Type                *bool       `json:"type"`
	RunUp               *string     `json:"run_up"`
	Environment         *string     `json:"environment"`
}

func GetFileDataExtraction(c *fiber.Ctx) error {
	apiKey := os.Getenv("GEMINI_API_KEY")

	file, err := c.FormFile("file")
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to get file")
	}
	ext := filepath.Ext(file.Filename) // Get extension (.png, .jpg, etc.)
	mimeType := mime.TypeByExtension(ext)

	ObjectType := c.FormValue("type", "billboard")
	isMultiple := c.FormValue("multiple", "no")

	// If MIME type is missing, try detecting from the extension
	if mimeType == "" {
		mimeType = file.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "application/octet-stream" // Default if unknown
		}
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to open file")
	}
	defer f.Close()

	// Read file data
	fileData, err := io.ReadAll(f)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to read file")
	}

	// Convert image to Base64
	encodedImage := base64.StdEncoding.EncodeToString(fileData)

	analyzeStatement := `one of the billboards or signages that is focused in this image and output a JSON array. For the object includes the following fields:`
	if isMultiple == "yes" {
		analyzeStatement = ` each billboards or signages contained in this image as a separate object and output a JSON array. For each object includes the following fields:`
	}

	instructions2 := fmt.Sprintf(`
    Analyze %s
        campaign_brand: Company advertising
        campaign_description: Brief ad description
        campaign_contacts: {phone: [integers], email: [strings]}
        campaign_site_url: [strings]
        location: Location if visible, else null
        billboard_measurements: {height: number, width: number, units: "meters"}
        target_audience: Likely audience
        target_age: "children", "youth", "adults", or "general"
        target_gender: "female", "male", "general"
        campaign_socials: {platform: handle} or null
        other_details: [{key: string, value: any, currency: string}] or []
        environment: "cluttered" or "solus"
        object_type: "billboard" or "signage"
        billboard_type: "Free Standing Billboard", "Free Standing Double Decker", "Wall Wrap", "Skysign", "Street Banner", "Bridge - Overpass", "Bridge - Footbridge", "Bridge - Railway", "Gantry", "Cantilever", "Free Standing Digital", "Wall Mounted Digital", "Wall Mounted INDOOR Digital LED", "Free Standing INDOOR Digital LED", "Wall Mounted INDOOR Digital LCD", "Free Standing Digital Indoor LCD", "Hanging INDOOR LCD", "Road Banner", "Wall Branding", "Bus Shelters", etc.
        structure: (for billboards: Bridge,digital, free standing, Gantry,hoarding,Hooding,Right,Sky, sky sign, wall wrap),(for branding: free standing banners, fish tank, backlit signs, stickers, etc.), (for signage: Name plate, Shop fascia, building fascia, 2D signage, 3D Signage, under canopy, back lit outdoor, etc.)
        material: "backlit", "digital", "flex", "LED", etc.
        illumination: "front" or "none"
        visibility: "Average", "Excellent", "Good", "Poor"
        run_up: ">100", "90-100", (>100, 90-100, 80-90, 70-80, 60-70, 50-60, 40-60, 30-40, 20-30, <20).
        angle: "double decker", "Head On", "Left", "Right", or null
        owner: {name, phone: [integers], email: [strings], address, website}
        type: true if %s is billboard_type, false otherwise
        percentage_accuracy: float
        additional_notes: Observations about the billboard"
    `, analyzeStatement, ObjectType)

	// Prepare the request payload
	requestBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: instructions2}, // User's instruction
					{InlineData: &Image{MimeType: mimeType, Data: encodedImage}}, // Image
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode request"})
	}

	// Send API request to Gemini
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to send API request")
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to read API response")
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to parse API response")
	}

	// Extract response text
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		//log the response
		fmt.Println(string(body))
		return utils.WriteError(c, fiber.StatusBadGateway, "No valid content received from Gemini API")
	}

	rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text

	// Extract JSON from the response text (clean up markdown formatting)
	re := regexp.MustCompile("```json\\s*([\\s\\S]*?)\\s*```")
	matches := re.FindStringSubmatch(rawJSON)

	var cleanedJSON string
	if len(matches) > 1 {
		cleanedJSON = matches[1] // Extract the JSON inside the Markdown block
	} else {
		// If no markdown code blocks, try to find the JSON directly
		cleanedJSON = strings.TrimSpace(rawJSON)
	}

	// Unmarshal extracted JSON into array of CampaignDetails structs
	var campaigns []CampaignDetails
	if err := json.Unmarshal([]byte(cleanedJSON), &campaigns); err != nil {
		// Log error details for debugging
		fmt.Println("JSON parsing error:", err)
		fmt.Println("Cleaned JSON:", cleanedJSON)
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to parse extracted JSON array")
	}

	return c.Status(fiber.StatusOK).JSON(campaigns)
}
