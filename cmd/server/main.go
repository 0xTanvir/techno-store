package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"techno-store/config"
	_ "techno-store/docs"
	"techno-store/internal/api/web"
	"techno-store/internal/infrastructure/datastores/pg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           technoStore API
// @version         1.0
// @description     Describes technoStore REST API
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  ./
// @query.collection.format multi

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load the config
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.Print()
	appConfig, err := config.Parse()
	if err != nil {
		slog.Error("Error parsing config", "cause", err)
	}

	// Get a datastore instance
	ds := pg.GetInstance(appConfig.Db)

	apiService := web.NewAPIService(*appConfig.Server, ds)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	apiService.InstallRoutes(router)

	// attaching swag to gin router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":" + appConfig.Server.Port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 2 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", "cause", err)
	}
	// catching ctx.Done(). timeout of 2 seconds.
	slog.Info("timeout of 2 seconds.")
	<-ctx.Done()

	slog.Info("Server exiting")
}

func initLogger() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
