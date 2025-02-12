package api

import (
	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	auth "bbscout/services/authorization"
	files "bbscout/services/files"
)

func SecuredRoutes(r fiber.Router) {
	i := r.Group("sl", middleware.CheckAccountAuthentication)

	// files
	i.Post("/upload/files", files.UploadFile)
	i.Get("/file/:fileName", files.GetFileByName)
	i.Get("/files", files.GetFiles)

	// change password
	i.Post("/change/password", auth.ChangePassword)

}
