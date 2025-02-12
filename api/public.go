package api

import (
	"github.com/gofiber/fiber/v2"

	services "bbscout/services/authorization"
)

func PublicRoutes(app fiber.Router) {
	app.Post("/login", services.Login)
	app.Get("/health", services.HealthCheck)

	//google verify Oauth2
	app.Post("/auth/google/verify", services.AuthGoogleVerify)

}
