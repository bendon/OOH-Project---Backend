package api

import (
	"github.com/gofiber/fiber/v2"

	services "bbscout/services/authorization"
)

func AuthorizedRoutes(r fiber.Router) {
	r.Get("/profile", services.GetUserProfile)
	r.Get("/accounts", services.GetUserAccounts)
	r.Post("/switch/account", services.PostSwtichAccounts)
}
