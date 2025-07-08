package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/roksky/bootstrap-api/config"
	"github.com/roksky/bootstrap-api/controller"
	"github.com/roksky/bootstrap-api/database"
	"github.com/roksky/bootstrap-api/helper"
	"github.com/roksky/bootstrap-api/job"
	"github.com/roksky/bootstrap-api/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type StratUpConfig struct {
	IntrospectURL string `json:"introspect_url"`
	SentryDSN     string `json:"sentry_dsn"`
}

// IntrospectURL := "https://auth.example.com/oauth/introspect"
// SentryDSN := "https://4edd0d587e0fdee6b69867ab46c2cb78@o4507588780425216.ingest.us.sentry.io/4507588782260224"
func StartApp(startupConfig StratUpConfig, provider controller.Provider) {
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

	helper.ErrorPanic(err)

	err = myDb.AutoMigrate()
	helper.ErrorPanic(err)

	// Router
	routeHandler, err := router.NewRouteHandler("/api")
	helper.ErrorPanic(err)

	// Enable Sentry
	routeHandler.EnableSentry(startupConfig.SentryDSN)

	// Enable CORS
	routeHandler.AllowCORS()

	err = routeHandler.EnableAuth(startupConfig.IntrospectURL, config.EnvConfigs.Auth.ClientId, config.EnvConfigs.Auth.ClientSecret)
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
