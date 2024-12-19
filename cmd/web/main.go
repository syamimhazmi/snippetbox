package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses
type Application struct {
	logger *slog.Logger
	db     *pgx.Conn
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

	dbDSN := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		os.Getenv("DB_CONNECTION"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	conn, err := pgx.Connect(context.Background(), dbDSN)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		logger.Error("Failed to established connection", "error", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	app := &Application{
		logger: logger,
		db:     conn,
	}

	mux := app.routes()

	logger.Info("starting server", "port", os.Getenv("APP_PORT"))

	err = http.ListenAndServe(os.Getenv("APP_PORT"), mux)
	logger.Error(err.Error())
	os.Exit(1)
}
