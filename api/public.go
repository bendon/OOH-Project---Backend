package api

import (
	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	services "bbscout/services/authorization"
	files "bbscout/services/files"
	gemini "bbscout/services/gemini"
)

func PublicRoutes(app fiber.Router) {
	app.Post("/auth/login", services.Login)
	app.Post("/auth/register", services.Register)
	app.Get("/health", services.HealthCheck)

	//google verify Oauth2
	app.Post("/auth/google/verify", services.AuthGoogleVerify)
	app.Get("/auth/file/:fileName", files.GetFileByName)
	app.Post("/auth/gemini/data/extraction", gemini.GetFileDataExtraction)

	auth := app.Group("auth", middleware.CheckAccountRefreshTokenAuthentication)
	auth.Post("/refresh/account", services.RefreshToken)

}
