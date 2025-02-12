package api

import (
	"github.com/gofiber/fiber/v2"

	services "bbscout/services/authorization"
)

func PublicRoutes(app fiber.Router) {
	app.Post("/login", services.Login)

}
