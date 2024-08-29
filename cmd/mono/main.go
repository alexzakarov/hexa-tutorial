package main

import (
	"context"
	"log"
	"main/config"
	_ "main/docs"
	coreServices "main/internal/v1/core_api/application/services"
	coreHandler "main/internal/v1/core_api/handler/http"
	coreRepos "main/internal/v1/core_api/infrastructure/repository"
	"main/pkg/databases/postgresql"
	"main/pkg/logger"
	"main/pkg/server"
	"main/pkg/utils/graceful_exit"
	"os"
)

// @title Core API
// @version 1.0
// @description Core service broker with REST endpoints
// @contact.email sercan@teamlify.co
// @BasePath /core/v1
func main() {
	os.Setenv("APP_ENV", "dev")

	cfg, errConfig := config.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", cfg.Sqs.AWS_ACCESS_KEY_ID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.Sqs.AWS_SECRET_ACCESS_KEY)

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Clients
	postgresqlDB, err := postgresql.NewPostgresqlDB(cfg)
	if err != nil {
		appLogger.Fatal("Error when tyring to connect to Postgresql")
	} else {
		appLogger.Info("Postgresql connected")
	}

	// Init repositories
	pgRepoCore := coreRepos.NewPostgresqlRepository(ctx, postgresqlDB)

	// Init services
	svcCore := coreServices.NewCoreService(cfg, pgRepoCore, appLogger)

	// Interceptors
	//

	servers := server.NewServer(cfg, &ctx, appLogger)

	httpServer, errHttpServer := servers.NewHttpServer()
	if errHttpServer != nil {
		println(errHttpServer.Error())
	}
	versioning := httpServer.Group(cfg.Server.API_VER)

	// Init handlers for HTTP Server
	hndlCore := coreHandler.NewHttpHandler(ctx, cfg, svcCore, appLogger)

	// Init routes for HTTP Server
	coreHandler.MapRoutes(cfg, hndlCore, versioning)

	servers.Listen(httpServer)

	// Exit from application gracefully
	graceful_exit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly")
}
