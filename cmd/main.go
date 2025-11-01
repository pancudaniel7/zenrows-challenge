package main

import (
	"sync"
	"zenrows-challenge/internal/adapter/http"
	"zenrows-challenge/internal/adapter/repo"
	"zenrows-challenge/internal/core/port"
	"zenrows-challenge/internal/core/usecase"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"gorm.io/gorm"

	"zenrows-challenge/internal/infra"
	"zenrows-challenge/internal/pkg/applog"
	"zenrows-challenge/internal/pkg/middleware"

	"github.com/go-playground/validator/v10"
)

var (
	logger applog.AppLogger
	v      *validator.Validate
	server *fiber.App
	db     *gorm.DB

	wg sync.WaitGroup

	// repo
	userRepo            port.UserRepo
	deviceTemplatesRepo port.DeviceTemplateRepo
	deviceProfileRepo   port.DeviceProfileRepo

	// service
	userSvc           port.AuthenticationService
	deviceTemplateSvc port.DeviceTemplateService
	deviceProfileSvc  port.DeviceProfileService

	// http handler
	deviceTemplateHandler port.DeviceTemplateHandler
	deviceProfileHandler  port.DeviceProfileHandler
)

func initComponents() {
	userRepo = repo.NewUserRepoImpl(logger, db)
	deviceTemplatesRepo = repo.NewDeviceTemplateRepoImpl(logger, db)
	deviceProfileRepo = repo.NewDeviceProfileRepoImpl(logger, db)

	userSvc = usecase.NewAuthenticationService(logger, userRepo)

	deviceTemplateSvc = usecase.NewDeviceTemplateServiceImpl(logger, deviceTemplatesRepo)
	deviceProfileSvc = usecase.NewDeviceProfileServiceImpl(logger, deviceProfileRepo, deviceTemplatesRepo, v)
	deviceTemplateHandler = http.NewDeviceTemplateHandlerImpl(logger, deviceTemplateSvc)
}

func initRoutes(server *fiber.App) {
	// Unprotected route
	server.Get("/health", func(c fiber.Ctx) error { return c.SendString("UP!") })

	// Protected routes group: apply BasicAuth to everything except /health
	protected := server.Group("/", middleware.BasicAuthCheckMiddleware(userSvc, v))
	protected.Get("/device-templates", deviceTemplateHandler.List)
}

func main() {
	if err := infra.LoadConfig(); err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	logger = applog.NewAppDefaultLogger()
	db = infra.ConnectToDatabase()
	v = validator.New()

	initComponents()
	server = infra.StartServer(logger, &wg)
	initRoutes(server)

	infra.GracefulShutdownServer(logger, &wg, server, nil)
}
