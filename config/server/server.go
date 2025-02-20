package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"bbscout/api"
	"bbscout/config/initializer"
	"bbscout/config/migration"
	"bbscout/middleware"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		addr: addr,
	}
}

func (s *ApiServer) Run() error {
	migration.InitializeMigrations()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://scout.edgetech.co.ke", // Replace with your frontend URL
		AllowMethods:     "PATCH,GET,POST,PUT,DELETE,PUT",
		AllowHeaders:     "Origin, Content-Type, Accept,Authorization",
		AllowCredentials: false,
	}))
	group := app.Group("/api/v1")

	api.PublicRoutes(group)

	authGroup := group.Group("/en", middleware.CheckAuthentication)

	api.AuthorizedRoutes(authGroup)
	api.SecuredRoutes(authGroup)
	app.Use(middleware.NotFoundMiddleware)

	// initialize operation account
	initializer.InitializerOperationAccount()

	fmt.Println("Starting server....")
	return app.Listen(s.addr)

}
