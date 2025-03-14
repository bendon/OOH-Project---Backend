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
}

func GetFileDataExtraction(c *fiber.Ctx) error {
	apiKey := os.Getenv("GEMINI_API_KEY")

	file, err := c.FormFile("file")
	if err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to get file")
	}
	ext := filepath.Ext(file.Filename) // Get extension (.png, .jpg, etc.)
	mimeType := mime.TypeByExtension(ext)

	ObjectType := c.FormValue("type")

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

	instructions := fmt.Sprintf(`Analyze the provided image of a billboard and extract the following information as a JSON object:
					campaign_brand: The brand or company advertising on the billboard.
					campaign_description: A brief description of the advertisement or promotion.
					campaign_contacts: * campaign_phone: The advert contact info as array integer default to empty array [].
					campaign_email: The email address for the campain as array string default to empty array [].
					location: The location of the billboard, if discernible from the image or context. If not visible, set to null.
					billboard_measurements:
					height: The estimated height of the billboard in meters. default 0.
					width: The estimated width of the billboard in meters. default  0.
					units: "meters"
					target_audience: A brief description of the likely target audience based on the ad content and billboard placement.
					Additional Notes as additional_notes:
					If the exact height and width cannot be determined, provide estimated values based on visual analysis and perspective.
					Include any observations or insights about the image, such as the type of billboard (digital, static), its surroundings, and the overall message of the advertisement. add also the percentage_accuracy for extraction as float. 
					campaign_site_url : Identify url on the image as site_url in array string if not place empty array.
					Extract target age as target_age either (children, youth,adults,general).
					Extract target gender as target_gender (female,male,general)
					Extract socials on the image as campaign_socials object as key and value e.g facebook,instagram,twitter,twitter or x ,linkedIn, github, WhatsApp etc as object string else empty null.
					Extract other details as other_details array object as key and value e.g [{key: price,value:100,currency: dollars }] etc as array string else empty array.
					Identify object in the image as either a billboard or signage as object_type as string else empty null.
					billboard_type: as either "Static Billboard", "Digital Billboard", "Banner Ads", "Wallscapes", "Mobile Billboards","Lamp Posts","Interactive Billboards" or null.
					structure : as either Bridge,digital, free standing, Gantry,hoarding,Hooding,Right,Sky, sky sign, wall wrap or null.
					material: as either backlit,digital,flex,LED,Vinyl,Sticker, Metal,Mesh or null.
					illumization: as either front or  none.
					visibility:  as either  Average, Excellent,Good,Poor.
					Angle: as either double decker, Head On,Left,Right or null.
					owner: the owner of the structure not advert. get the owner_name, owner_phone as array of integer, owner_email array of string, owner_address, owner_website as object string else empty null. 
					type: as either true or false if the object type is %s.
					Format the output as a JSON object with the specified fields.`, ObjectType)

	// Prepare the request payload
	requestBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: instructions}, // User's instruction
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

	// Remove Markdown code block (```json ... ```)
	re := regexp.MustCompile("(?s)```json\\n(.*?)\\n```")
	matches := re.FindStringSubmatch(rawJSON)

	var cleanedJSON string
	if len(matches) > 1 {
		cleanedJSON = matches[1] // Extract the JSON inside the Markdown block
	} else {
		cleanedJSON = strings.TrimSpace(rawJSON) // If no Markdown, use the raw text
	}

	// Unmarshal extracted JSON into CampaignDetails struct
	var campaign CampaignDetails
	if err := json.Unmarshal([]byte(cleanedJSON), &campaign); err != nil {
		return utils.WriteError(c, fiber.StatusBadGateway, "Failed to parse extracted JSON")
	}

	return c.Status(fiber.StatusOK).JSON(campaign)

}
