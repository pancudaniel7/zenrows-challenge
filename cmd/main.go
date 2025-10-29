package main

import (
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"zenrows-challenge/internal/infra"
	"zenrows-challenge/internal/pkg/applog"
)

var (
	logger applog.AppLogger
	server *fiber.App

	wg sync.WaitGroup
)

func main() {
	if err := infra.LoadConfig(); err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	logger = applog.NewAppDefaultLogger()
	server = infra.StartServer(logger, &wg)

	//scb := func() error {
	//	logger.Debug("Executing graceful shutdown callback...")
	//	return nil
	//}
	infra.ShutdownServer(logger, &wg, server, nil)

}
