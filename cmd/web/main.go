package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/syamimhazmi/snippetbox/internal/database"
	"github.com/syamimhazmi/snippetbox/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses
type Application struct {
	env      string
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// Use the slow.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	conn := database.New(logger)

	defer conn.Close(context.Background())

	app := &Application{
		env:      os.Getenv("APP_ENV"),
		logger:   logger,
		snippets: &models.SnippetModel{DB: conn},
	}

	mux := app.routes()

	logger.Info("starting server", "port", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(os.Getenv("APP_PORT"), mux)
	logger.Error(err.Error())
	os.Exit(1)
}
