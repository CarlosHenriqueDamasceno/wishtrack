package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/suggestion"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	httputils "github.com/CarlosHenriqueDamasceno/wishtrack/pkg/http"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func run(conf *server.Config, logger *slog.Logger) error {
	conn, err := database.New(
		conf.Database.Dsn,
		conf.Database.MaxOpenConns,
		conf.Database.MaxIdleConns,
		conf.Database.MaxIdleTime,
	)
	if err != nil {
		return err
	}

	conf.SetDatabaseConnection(conn)

	logger.Info("Database connection established.")

	userRepository := user.NewDatabaseRepository(conf.Database.Conn)
	contentRepository := content.NewDatabaseRepository(conf.Database.Conn)
	jwtAuth := user.NewJwtAuthenticator(conf.Auth.Key, conf.Auth.Iss, conf.Auth.Aud, conf.Auth.Exp)
	userService := user.NewService(userRepository, jwtAuth)
	contentService := content.NewService(contentRepository)

	client := &http.Client{
		Transport: httputils.NewAuthenticatedTransport(conf.Providers.TMDB.ApiKey),
	}

	tmdbSuggester, err := suggestion.NewSuggester(
		suggestion.TMDB,
		client,
		conf.Providers.TMDB.BaseUrl,
		contentRepository,
	)
	if err != nil {
		logger.Error("Error creating TMDB suggester", "error", err)
	}

	personalSuggester, err := suggestion.NewSuggester(
		suggestion.PERSONAL,
		client,
		conf.Providers.TMDB.BaseUrl,
		contentRepository,
	)
	if err != nil {
		logger.Error("Error creating TMDB suggester", "error", err)
	}

	suggesters := []suggestion.Suggester{tmdbSuggester, personalSuggester}
	suggestionService := suggestion.NewService(suggesters)

	api := server.NewApi(http.NewServeMux(), conf, logger, userService, contentService, suggestionService)
	server := http.Server{
		Addr:         conf.Address,
		Handler:      api,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Info("signal caught", "signal", s.String())

		shutdown <- server.Shutdown(ctx)
	}()

	logger.Info("server has started", "addr", conf.Address, "env", conf.Env)

	err = server.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	logger.Info("server has stopped", "addr", conf.Address, "env", conf.Env)

	return nil
}

// @title						WishTrack API
// @version					1.0
// @description				A app to manage all the things you wanna see (or had saw).
// @contact.name				API Support
// @contact.email				carlos@wishtrack.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8080
// @BasePath					/api/v1
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	config := server.LoadEnv()

	if err := run(config, logger); err != nil {
		logger.Error("Server error", "error", err.Error())
		os.Exit(1)
	}
}
