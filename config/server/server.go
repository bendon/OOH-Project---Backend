package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"bbscout/api"
	"bbscout/config/initializer"
	"bbscout/config/migration"
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
		AllowOrigins:     "https://inventory.diracks.com, http://localhost:5173", // Replace with your frontend URL
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,PUT",
		AllowHeaders:     "Origin, Content-Type, Accept,Authorization",
		AllowCredentials: true,
	}))
	group := app.Group("/api/v1")

	api.PublicRoutes(group)


	// initialize operation account
	initializer.InitializerOperationAccount()

	fmt.Println("Starting server....")
	return app.Listen(s.addr)

}
