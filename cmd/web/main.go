package main

import (
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses
type Application struct {
	logger *slog.Logger
}

func main() {
	// Use the slow.New() function to initialize a new structured logger, which
	// writes to the standard out stream and uses the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	_, err := os.Stat(".env")
	if err != nil {
		logger.Error(".env file not found")
		os.Exit(1)
	}

	app := &Application{
		logger: logger,
	}

	mux := app.routes()

	logger.Info("starting server on :4000")

	err = http.ListenAndServe(":4000", mux)
	logger.Error(err.Error())
	os.Exit(1)
}
