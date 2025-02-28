package api

import (
	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	services "bbscout/services/authorization"
	files "bbscout/services/files"
)

func PublicRoutes(app fiber.Router) {
	app.Post("/auth/login", services.Login)
	app.Get("/health", services.HealthCheck)

	//google verify Oauth2
	app.Post("/auth/google/verify", services.AuthGoogleVerify)
	app.Get("/file/:fileName", files.GetFileByName)

	auth := app.Group("auth", middleware.CheckAccountRefreshTokenAuthentication)
	auth.Post("/refresh/account", services.RefreshToken)

}
