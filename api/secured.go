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
	i.Get("/staff/:staffId", organization.GetOrganizationStaffById)
	i.Get("/organization/user/analytics", organization.OrganizationUserAnalytics)
	i.Get("/permission/staff/:staffId", organization.GetStaffPermissions)
	i.Post("/permisions/staff/:staffId/update", organization.UpdateStaffPermissions)

	//roles
	i.Get("/roles", organization.GetOrganizationRoles)
	i.Post("/role", organization.CreateOrganizationRoles)
	i.Post("/role/update", organization.UpdateOrganizationRoles)
	i.Get("/permissions", organization.GetPermissions)

	// billboard
	i.Post("/billboard", organization.CreateBillBoard)
	i.Get("/billboard/:id", organization.GetBillBoardById)
	i.Get("/billboards", organization.GetBillBoards)
	i.Delete("/delete/billboard/:id", organization.DeleteBillBoard)
	i.Put("/update/billboard/:id", organization.UpdateBillBoard)
	i.Get("/my/billboardds/uploads", organization.MyBillBoardsUploads)
	i.Get("/my/billboards/:billboardId", organization.MyBillBoardById)
	i.Get("/billboards-daily-filter", organization.GetBillboardDailyFilter)

	// billboard campaigns
	i.Post("/billboard/campaign", organization.CreateBillboardCampaign)
	i.Get("/billboards/organization/campaigns", organization.GetBillBaordsOrganizationCampaigns)
	i.Get("/billboard/campaign/:campaignId", organization.GetCampaignById)
	i.Get("/billboard/history/:billboardId/campaigns", organization.GetBillboardHistryCampaigns)
	i.Put("/campaign/update/:campaignId", organization.UpdateBillboardCampaign)
	i.Delete("/campaign/delete/:campaignId", organization.DeleteBillboardCampaign)
	i.Post("/campaign/:campaignId/close", organization.CloseBillboardCampaign)
	i.Post("/create/billboard/types", organization.CreateBillboardTypes)
	i.Delete("/delete/billboard/type/:typeId", organization.DeleteBillboardType)
	i.Get("/list/billboard/types", organization.GetBillboardTypes)
	i.Get("/realted/billboards/:billboardCode", organization.RelatedBillBoardByCode)

	// send email
	i.Post("/send/email", organization.SendEmail)

	// reports
	i.Get("/report/billboard/types", organization.BillboardTypeReports)
	i.Get("/report/billboard/organization", organization.BillboardOrganizationReports)
	i.Get("/report/billboard/locations", organization.BillboardLocationReports)
	i.Get("/report/billboard/weekly", organization.BillboardWeeklyReports)
	i.Get("/report/billboard/monthly", organization.BillboardMonthlyReports)
	i.Get("/report/billboard/yearly", organization.BillboardYearlyReports)
	i.Get("/report/billboard/user/uploads", organization.GetMyUploadsSummary)
	i.Get("/report/billboard/user/organization", organization.GetUserOrganizationUploadsSummary)
	i.Get("/report/billboard/user/weekly", organization.GetUserUploadReportsWeekly)
	i.Get("/report/billboard/user/monthly", organization.GetUserUploadReportsMonthly)
	i.Get("/report/billboard/user/yearly", organization.GetUserUploadReportsYearly)
}
