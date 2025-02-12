package api

import (
	"github.com/gofiber/fiber/v2"

	services "bbscout/services/authorization"
)

func PublicRoutes(app fiber.Router) {
	app.Post("/login", services.Login)
	app.Post("/upload/files", services.UploadFile)
	app.Get("/file/:fileName", services.GetFileByName)
	app.Get("/files", services.GetFiles)

}
