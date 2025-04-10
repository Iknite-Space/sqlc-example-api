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
	ctx := context.Background() //creates a route context for all DB operations
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
	db, err := pgxpool.New(ctx, dbConnectionURL) //connect to postgresSQL using pgxpol (reusable)
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
	handler := api.NewMessageHandler(querier).WireHttpHandler()

	// And finally we start the HTTP server on the configured port.
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	return nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(cfg *Config) error {
	if _, err := os.Stat(".env"); err == nil { //checks for the presense of an .env
		err = godotenv.Load() //loads vars from .env to os.Environ, so that we can access them with os.Getenv("DB_USER")
		if err != nil {
			return fmt.Errorf("failed to load env file: %w", err)
		}
	}
	//conf.Parse() uses reflection to autofill the struct from env vars
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
	queryValues := url.Values{} //a map for url query string (e.g sslMode=require)
	if config.TLSDisabled {
		queryValues.Add("sslmode", "disable")
	} else {
		queryValues.Add("sslmode", "require")
	} //decide whether to disable TLS or not

	//construct the actual URL
	dbURL := url.URL{
		Scheme:   "postgres",                                         //this makes it makes it postgres:// URL
		User:     url.UserPassword(config.DBUser, config.DBPassword), //injects the username and password
		Host:     fmt.Sprintf("%s:%d", config.DBHost, config.DBPort), //localhost:8085
		Path:     config.DBName,                                      // this is the name of the database
		RawQuery: queryValues.Encode(),                               //adds the  ?sslmode=require part
	}

	return dbURL.String() //and finally returns somethings like this postgres://ichami:supersecret@localhost:5432/mydb?sslmode=disable
}
