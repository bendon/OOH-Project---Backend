package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"bbscout/models"
)

var JwtKey = []byte("jOvJrciv7VGxFDy1KVGRg3WNiv2KngQa5")

type Claims struct {
	OwnerID  uuid.UUID `json:"accessor"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type AccountClaims struct {
	OwnerID     uuid.UUID  `json:"accessor"`
	Accessing   uuid.UUID  `json:"accessing"`
	Accessor    *uuid.UUID `json:"access"`
	Permissions []string   `json:"permissions"`
	Username    string     `json:"username"`
	jwt.RegisteredClaims
}

type AccountBranchClaimResponse struct {
	OwnerID     uuid.UUID `json:"accessor"`
	Accessing   uuid.UUID `json:"accessing"`
	Accessor    uuid.UUID `json:"access"`
	Username    string    `json:"username"`
	Permissions []string  `json:"permissions"`
}

type LocationData struct {
	Country    string  `json:"country_name"`
	Region     string  `json:"region_name"`
	City       string  `json:"city"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	TimeZone   string  `json:"timezone"`
	RegionCode string  `json:"region_code"`
	PublicIp   string  `json:"ip"`
}

type TokenResponse struct {
	User         models.UserModel         `json:"user"`
	Account      *models.UserAccountModel `json:"account,omitempty"`
	AccessToken  string                   `json:"accessToken"`
	RefreshToken string                   `json:"refreshToken"`
	Permissions  *[]string                `json:"permissions,omitempty"`
}

type contextKey string

const UserIDKey contextKey = "userId"
