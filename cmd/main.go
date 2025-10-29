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
)

var (
	logger applog.AppLogger
	server *fiber.App
	db     *gorm.DB

	wg sync.WaitGroup

	// repo
	userRepo            port.UserRepo
	deviceTemplatesRepo port.DeviceTemplateRepo

	// service
	userSvc           port.AuthenticationService
	deviceTemplateSvc port.DeviceTemplateService

	// http handler
	deviceTemplateHandler port.DeviceTemplateHandler
)

func main() {
	if err := infra.LoadConfig(); err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	logger = applog.NewAppDefaultLogger()

	db = infra.ConnectToDatabase()
	server = infra.StartServer(logger, &wg)

	userRepo = repo.NewUserRepoImpl(logger, db)
	deviceTemplatesRepo = repo.NewDeviceTemplateRepoImpl(logger, db)

	userSvc = usecase.NewAuthenticationService(logger, &wg, userRepo)
	deviceTemplateSvc = usecase.NewDeviceTemplateServiceImpl(logger, deviceTemplatesRepo)

	deviceTemplateHandler = http.NewDeviceTemplateHandlerImpl(logger, deviceTemplateSvc)

	//scb := func() error {
	//	logger.Debug("Executing graceful shutdown callback...")
	//	return nil
	//}
	infra.ShutdownServer(logger, &wg, server, nil)

}
