package api

import (
	"github.com/gofiber/fiber/v2"

	"bbscout/middleware"
	auth "bbscout/services/authorization"
	files "bbscout/services/files"
	organization "bbscout/services/organization"
)

func SecuredRoutes(r fiber.Router) {
	i := r.Group("sl", middleware.CheckAccountAuthentication)

	// files
	i.Post("/upload/files", files.UploadFile)
	i.Get("/file/:fileName", files.GetFileByName)
	i.Get("/files", files.GetFiles)

	// change password
	i.Post("/change/password", auth.ChangePassword)

	// Organization
	i.Get("/organization/profile", organization.GetOrganizationProfile)
	i.Post("/new/staff", organization.NewOrganizationStaff)
	i.Get("/staffs", organization.OrganizationStaffs)
	i.Get("/organization/user/analytics", organization.OrganizationUserAnalytics)

	// billboard
	i.Post("/billboard", organization.CreateBillBoard)
	i.Get("/billboard/:id", organization.GetBillBoardById)
	i.Get("/billboards", organization.GetBillBoards)
	i.Delete("/delete/billboard/:id", organization.DeleteBillBoard)
	i.Put("/update/billboard/:id", organization.UpdateBillBoard)

}
