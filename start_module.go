package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/roksky/bootstrap-api/auth"
	"github.com/roksky/bootstrap-api/config"
	"github.com/roksky/bootstrap-api/controller"
	"github.com/roksky/bootstrap-api/database"
	"github.com/roksky/bootstrap-api/helper"
	"github.com/roksky/bootstrap-api/job"
	"github.com/roksky/bootstrap-api/router"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func startApp() {
	file, err := os.OpenFile("gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Msgf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Configure Zerolog to write to the file
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339})

	log.Info().Msg("Started Server!")

	config.InitEnvConfigs()

	db := config.DatabaseConnection()
	myDb := database.NewDatabase(db)
	validate := validator.New()

	authDB := auth.NewAuthDatabase(db)
	err = authDB.AutoMigrate()
	helper.ErrorPanic(err)

	err = myDb.AutoMigrate()
	helper.ErrorPanic(err)

	// controllers
	provider := controller.NewProvider(db, validate)

	// Router
	routeHandler, err := router.NewRouteHandler("/api", db, authDB)
	helper.ErrorPanic(err)

	// Enable Sentry
	routeHandler.EnableSentry()

	// Enable CORS
	routeHandler.AllowCORS()

	err = routeHandler.EnableAuth(config.EnvConfigs.Auth.ClientId, config.EnvConfigs.Auth.ClientSecret)
	helper.ErrorPanic(err)

	routeHandler.RegisterRoutes(provider.GetControllers())
	initJobExecutor(provider.GetJobs())

	httpServer := &http.Server{
		Addr:    ":" + config.EnvConfigs.ServerPort,
		Handler: routeHandler.GetEngine(),
	}

	err = httpServer.ListenAndServe()
	helper.ErrorPanic(err)
}

func initJobExecutor(jobs []job.Job) {
	je := job.NewJobExecutor()

	for _, myJob := range jobs {
		je.RegisterJob(myJob)
	}
	je.Cron.Start()
	fmt.Println("Job run log:", je.JobRunLog)
}
