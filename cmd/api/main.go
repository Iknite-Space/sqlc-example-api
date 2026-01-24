// Package main implements the entry point for the application.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Iknite-Space/sqlc-example-api/api"
	"github.com/Iknite-Space/sqlc-example-api/db/repo"
)

// DBConfig holds the database configuration. This struct is populated from the .env in the current directory.
type DBConfig struct {
	DBUser      string `conf:"env:DB_USER,required"`
	DBPassword  string `conf:"env:DB_PASSWORD,required,mask"`
	DBHost      string `conf:"env:DB_HOST,required"`
	DBPort      uint16 `conf:"env:DB_PORT,required"`
	DBName      string `conf:"env:DB_Name,required"`
	TLSDisabled bool   `conf:"env:DB_TLS_DISABLED"`
}

// Config holds the application configuration. This struct is populated from the .env in the current directory.
type Config struct {
	ListenPort     uint16 `conf:"env:LISTEN_PORT,required"`
	MigrationsPath string `conf:"env:MIGRATIONS_PATH,required"`
	DB             DBConfig
}

func main() {

	// We call run() here because main cannot return an error. If run() returns an error we print the error and exit.
	// This is a common pattern in Go applications to handle errors gracefully.
	err := run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// run initializes the application and starts the server.
// It loads the configuration, sets up the database connection, and starts the HTTP server.
func run() error {
	ctx := context.Background()
	config := Config{}

	// We load the configuration from the .env file in the current directory and populate the Config struct.
	// If the .env file is not found, or if any of the required configuration values are missing, an error is returned.
	err := LoadConfig(&config)
	if err != nil {
		fmt.Println("Error loading config:", err)
		fmt.Println("Have you configured your .env with the required variables?")
		return err
	}

	// We use the configuration values to get the database connection URL.
	dbConnectionURL := getPostgresConnectionURL(config.DB)
	db, err := pgxpool.New(ctx, dbConnectionURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// We use the database connection to run the migrations.
	// This will create or update all the required database tables.
	err = repo.Migrate(dbConnectionURL, config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	querier := repo.New(db)

	// We create a new http handler using the database querier.
	handler := api.NewAPIHandler(querier).WireHttpHandler()

	// And finally we start the HTTP server on the configured port.
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	return nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(cfg *Config) error {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}

	_, err := conf.Parse("", cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return err
		}

		return err
	}

	return nil
}

// getPostgresConnectionURL constructs the PostgreSQL connection URL from the provided configuration.
func getPostgresConnectionURL(config DBConfig) string {
	queryValues := url.Values{}
	if config.TLSDisabled {
		queryValues.Add("sslmode", "disable")
	} else {
		queryValues.Add("sslmode", "require")
	}

	dbURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.DBUser, config.DBPassword),
		Host:     fmt.Sprintf("%s:%d", config.DBHost, config.DBPort),
		Path:     config.DBName,
		RawQuery: queryValues.Encode(),
	}

	return dbURL.String()
}
